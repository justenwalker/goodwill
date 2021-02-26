// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task;

import com.walmartlabs.concord.ApiClient;
import com.walmartlabs.concord.client.ApiClientConfiguration;
import com.walmartlabs.concord.client.ApiClientFactory;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.*;
import tech.justen.concord.goodwill.service.*;

import java.io.*;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.SecureRandom;
import java.util.*;
import java.util.concurrent.ExecutorService;

import static java.lang.String.format;

public class TaskCommon {

    private static final Logger log = LoggerFactory.getLogger(TaskCommon.class);

    private static final Logger processLog = LoggerFactory.getLogger("processLog");

    private static final String MAGIC_VALUE = "d0c08ee0-a663-4a6b-ad5e-00a5fca1e5cf";

    private final TaskConfig config;

    private final TaskParams params;

    private final DependencyManager dependencyManager;

    private final DockerService dockerService;

    private final ContextService contextService;

    private final SecretService secretService;

    private final ExecutorService executor;

    private final LockService lockService;

    private final ApiClientConfiguration apiClientConfig;

    private final ApiClientFactory apiClientFactory;

    public TaskCommon(
            TaskConfig config,
            TaskParams params,
            DependencyManager dependencyManager,
            ContextService executionService,
            DockerService dockerService,
            SecretService secretService,
            LockService lockService,
            ExecutorService executor,
            ApiClientConfiguration apiClientConfig,
            ApiClientFactory apiClientFactory) {
        this.config = config;
        this.params = params;
        this.dependencyManager = dependencyManager;
        this.contextService = executionService;
        this.lockService = lockService;
        this.dockerService = dockerService;
        this.secretService = secretService;
        this.executor = executor;
        this.apiClientConfig = apiClientConfig;
        this.apiClientFactory = apiClientFactory;
    }

    public ApiClient getSessionClient() {
        return apiClientFactory.create(apiClientConfig);
    }

    public String compileInDocker() throws java.lang.Exception {
        log.info("Compiling goodwill binary in Docker");
        Path goodwillContainerBinPath = Paths.get("/workspace", params.getDirectory(), "goodwill");
        Path goodwillBinPath = Paths.get(config.workingDirectory().toString(), params.getDirectory(), "goodwill");
        File goodwillBin = goodwillBinPath.toFile();
        if (!goodwillBin.exists()) {
            String binClasspath = String.format("/go/%s/%s/goodwill%s", "linux", "amd64", "");
            try (InputStream link = (getClass().getResourceAsStream(binClasspath))) {
                Files.copy(link, goodwillBin.getAbsoluteFile().toPath());
            }
        }
        if (!goodwillBin.canExecute()) {
            goodwillBin.setExecutable(true);
        }
        Map<String, String> env = new HashMap<>();
        env.put("GOROOT", "/usr/local/go");
        Path out = Paths.get("/workspace", params.getFlowBinary());
        DockerContainer container = new DockerContainer();
        container.entryPoint = goodwillContainerBinPath.toString();
        container.command = Arrays.asList("-debug", "-dir", "/workspace", "-out", out.toString());
        container.image = "golang:1.16";
        container.env = env;
        container.workDir = "/workspace";
        container.debug = true;
        int result = dockerService.start(container, line -> {
            processLog.info("COMPILE: {}", line);
        }, line -> {
            processLog.info("COMPILE: {}", line);
        });
        if (result != 0) {
            throw new RuntimeException("goodwill exited unsuccessfully. See output logs for details.");
        }
        return Paths.get(config.workingDirectory().toString(), params.getFlowBinary()).toString();
    }

    public String compileOnHost() throws java.lang.Exception {
        log.info("Compiling goodwill binary on the agent");
        String os = params.getGoOS();
        String arch = params.getGoArch();
        File goodwillBin = params.getBinaryOutPath(config.workingDirectory()).toFile();
        if (!goodwillBin.exists()) {
            String binClasspath = params.getBinaryClasspath();
            try (InputStream link = (getClass().getResourceAsStream(binClasspath))) {
                Files.copy(link, goodwillBin.getAbsoluteFile().toPath());
            }
        }
        if (!goodwillBin.canExecute()) {
            if (!goodwillBin.setExecutable(true)) {
                throw new RuntimeException("Cannot make set executable bit");
            }
        }
        File out = new File(config.workingDirectory().toString(), params.getFlowBinary());
        Map<String, String> env = new HashMap<>();
        if (params.installGo) {
            Path goRoot = installGo();
            String path = System.getenv("PATH");
            env.put("PATH", path + File.pathSeparatorChar + Paths.get(goRoot.toString(), "bin").toString());
            env.put("GOROOT", goRoot.toString());
        }
        List<String> cmd = new ArrayList<>();
        cmd.add(goodwillBin.toString());
        cmd.add("-os");
        cmd.add(os);
        cmd.add("-arch");
        cmd.add(arch);
        cmd.add("-debug");
        cmd.add("-dir");
        cmd.add(config.workingDirectory().toString());
        cmd.add("-out");
        cmd.add(out.toString());
        exec(env, cmd.toArray(new String[0]));
        return out.toString();
    }

    public void execute() throws java.lang.Exception {
        Path workDir = config.workingDirectory();
        File goodwillDir = new File(workDir.toString(), params.getDirectory());
        goodwillDir.mkdir();
        File goodwillBin = new File(workDir.toString(), params.getFlowBinary());
        String commandPath = goodwillBin.toString();
        if (!goodwillBin.exists()) {
            log.debug("Goodwill binary {} does not exist", goodwillBin.toString());
            if (params.useDockerImage) {
                commandPath = compileInDocker();
            } else {
                commandPath = compileOnHost();
            }
        }
        String taskName = params.getTask();
        goodwillBin = new File(commandPath);
        boolean v = goodwillBin.setExecutable(true);
        Server server = null;
        CertUtils.CA ca = CertUtils.generateCA();
        InputStream caCert = ca.getCACertInputStream();
        InputStream caKey = ca.getCAKeyInputStream();
        File caFile = new File(goodwillDir, "ca.crt");
        File certFile = new File(goodwillDir, "client.crt");
        File keyFile = new File(goodwillDir, "client.key");
        ca.generatePKI(caFile, certFile, keyFile);
        try {
            ApiClient apiClient = apiClientFactory.create(apiClientConfig);
            int port = 0;
            long sleepMillis = 0;
            IOException startException = null;
            for (int i = 0; i < 10; i++) {
                Thread.sleep(sleepMillis);
                try {
                    port = randomPort();
                    server = ServerBuilder.forPort(port)
                            .useTransportSecurity(caCert, caKey)
                            .addService(new GrpcDockerService(dockerService))
                            .addService(new GrpcConfigService(apiClientConfig, config))
                            .addService(new GrpcContextService(contextService))
                            .addService(new GrpcSecretService(config, secretService, apiClient))
                            .addService(new GrpcLockService(lockService))
                            .build();
                    server.start();
                    startException = null;
                } catch (IOException e) {
                    startException = e;
                    sleepMillis = i * 1000;
                    port = 0;
                    log.warn("GRPC service failed to start on port {}, trying again in {} ms", port, sleepMillis);
                    continue;
                }
                break;
            }
            if (startException != null) {
                log.error("GRPC Service failed to start,");
                throw startException;
            }
            Map<String, String> env = new HashMap<>();
            env.put("GRPC_ADDR", format(":%d", port));
            env.put("GRPC_MAGIC_KEY", MAGIC_VALUE);
            env.put("GRPC_CA_CERT_FILE", caFile.getAbsolutePath());
            env.put("GRPC_CLIENT_CERT_FILE", certFile.getAbsolutePath());
            env.put("GRPC_CLIENT_KEY_FILE", keyFile.getAbsolutePath());
            env.put("CONCORD_ORG_NAME", config.orgName());
            env.put("CONCORD_PROCESS_ID", config.processId());
            env.put("CONCORD_WORKING_DIRECTORY", config.workingDirectory().toString());
            exec(env, goodwillBin.getAbsolutePath(), taskName);
        } finally {
            if (server != null) {
                server.shutdown();
            }
        }
    }

    private int randomPort() {
        SecureRandom s = new SecureRandom();
        return 49152 + s.nextInt(65535 - 49152);
    }

    private void exec(Map<String, String> env, String... command) throws IOException, InterruptedException {
        Path workDir = config.workingDirectory();
        ProcessBuilder pb = new ProcessBuilder();
        pb.command(command);
        if (env != null) {
            for (Map.Entry<String, String> e : env.entrySet()) {
                pb.environment().put(e.getKey(), e.getValue());
            }
        }
        pb.directory(workDir.toFile());

        String commandString = String.join(" ", pb.command());
        log.debug("Exec Goodwill Task: [{}]", commandString);

        Process p = pb.start();
        executor.execute(() -> {
            try (BufferedReader stdout = new BufferedReader(new InputStreamReader(p.getInputStream()))) {
                String line;
                while ((line = stdout.readLine()) != null) {
                    processLog.info("[OUT] GOODWILL: {}", line);
                }
            } catch (IOException e) {
                log.error("error reading stdout", e);
            }
        });
        executor.execute(() -> {
            try (BufferedReader stderr = new BufferedReader(new InputStreamReader(p.getErrorStream()))) {
                String line;
                while ((line = stderr.readLine()) != null) {
                    processLog.info("[ERR] GOODWILL: {}", line);
                }
            } catch (IOException e) {
                log.error("error reading stderr", e);
            }
        });

        int rc = p.waitFor();

        if (rc != 0) {
            log.warn("call ['{}'] -> finished with code {}",
                    commandString, rc);
            throw new RuntimeException("goodwill command failed");
        }
    }

    private Path installGo() throws IOException {
        String version = params.getGoVersion();
        Path workDir = config.workingDirectory();
        Path goRootDir = Paths.get(workDir.toString(), params.getDirectory(), "go");
        File goInstall = Paths.get(workDir.toString(), params.getDirectory(), "go", "bin", "go").toFile();
        if (goInstall.canExecute()) {
            log.info("Go already installed at {}", goInstall.toString());
            return goRootDir;
        }
        log.info("Installing Go {}", version);
        Path tar = dependencyManager.resolve(URI.create(params.getGoDownloadURL(version)));
        TarUtils.extractTarball(tar, Paths.get(workDir.toString(), params.getDirectory()));
        log.info("Go installed at {}", goInstall.toString());
        return goRootDir;
    }

}
