// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task;

import org.apache.commons.lang3.SystemUtils;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Map;

public class TaskParams {

    public static final String GO_DOCKER_IMAGE_KEY = "goDockerImage";

    public static final String GOOS_KEY = "goos";

    public static final String GOARCH_KEY = "goarch";

    public static final String TASK_NAME_KEY = "task";

    public static final String BINARY_KEY = "binary";

    public static final String DIRECTORY_KEY = "dir";

    public static final String DEBUG_KEY = "debug";

    public static final String USE_DOCKER_IMAGE_KEY = "useDocker";

    public static final String INSTALL_GO_KEY = "installGo";

    public static final String GO_VERSION_KEY = "goVersion";

    public static final String GOPROXY_KEY = "GOPROXY";

    public static final String GONOPROXY_KEY = "GONOPROXY";

    public static final String GOPRIVATE_KEY = "GOPRIVATE";

    public static final String GOSUMDB_KEY = "GOSUMDB";

    public static final String GONOSUMDB_KEY = "GONOSUMDB";

    private static final String DEFAULT_GO_VERSION = "1.16.2";

    private static final String DEFAULT_GO_OS = "linux";

    private static final String DEFAULT_GO_ARCH = "amd64";

    private static final String DEFAULT_GOODWILL_DIR = ".goodwill";

    private static final String DEFAULT_GO_DOCKER_IMAGE = "golang:1.16";

    private static final String DEFAULT_BIN = "goodwill.flow";

    private static final String DEFAULT_TASK = "Default";

    public String goProxy;

    public String goNoProxy;

    public String goPrivate;

    public String goSumDB;

    public String goNoSumDB;

    public String goDockerImage;

    public String goDownloadURL;

    public String goOS;

    public String goArch;

    public String taskName;

    public String flowBinary;

    public String flowDirectory;

    public boolean useDockerImage;

    public boolean installGo;

    public boolean debug;

    public String goVersion;

    public TaskParams() {

    }

    public void setGoEnvironment(Map<String, String> env) {
        setEnv(env, "GOPRIVATE", goPrivate);
        setEnv(env, "GOPROXY", goProxy);
        setEnv(env, "GONOPROXY", goNoProxy);
        setEnv(env, "GOSUMDB", goSumDB);
        setEnv(env, "GONOSUMDB", goNoSumDB);
    }

    private void setEnv(Map<String, String> env, String key, String value) {
        if (value != null && !value.isEmpty()) {
            env.put(key, value);
            return;
        }
        value = System.getenv(key);
        if (value != null && !value.isEmpty()) {
            env.put(key, value);
        }
    }

    public String getGoDownloadURL(String version) {
        if (goDownloadURL != null && !goDownloadURL.isEmpty()) {
            return goDownloadURL;
        }
        return String.format("https://golang.org/dl/go%s.%s-%s.tar.gz", version, getGoOS(), getGoArch());
    }

    public String getTask() {
        if (taskName == null || taskName.isEmpty()) {
            return DEFAULT_TASK;
        }
        return taskName;
    }

    public String getFlowBinary() {
        if (flowBinary == null || flowBinary.isEmpty()) {
            return DEFAULT_BIN;
        }
        return flowBinary;
    }

    public String getDirectory() {
        if (flowDirectory == null || flowDirectory.isEmpty()) {
            return DEFAULT_GOODWILL_DIR;
        }
        return flowDirectory;
    }

    public Path getBinaryOutPath(Path workingDirectory) {
        String ext = "";
        if (getGoOS().equals("windows")) {
            ext = ".exe";
        }
        return Paths.get(workingDirectory.toString(), getDirectory(), String.format("goodwill%s", ext));
    }

    public String getBinaryClasspath() {
        String ext = "";
        String os = getGoOS();
        String arch = getGoArch();
        if (os.equals("windows")) {
            ext = ".exe";
        }
        return String.format("/go/goodwill_%s_%s%s", os, arch, ext);
    }

    public String getGoOS() {
        if (goOS == null || goOS.isEmpty()) {
            if (SystemUtils.IS_OS_MAC_OSX) {
                return "darwin";
            }
            if (SystemUtils.IS_OS_WINDOWS) {
                return "windows";
            }
            if (SystemUtils.IS_OS_LINUX) {
                return "linux";
            }
            return DEFAULT_GO_OS;
        }
        return goOS;
    }

    public String getGoArch() {
        if (goArch == null || goArch.isEmpty()) {
            switch (SystemUtils.OS_ARCH) {
                case "amd64":
                case "x86_64":
                    return "amd64";
                case "x86":
                case "i386":
                    return "386";

            }
            return DEFAULT_GO_ARCH;
        }
        return goArch;
    }

    public String getGoDockerImage() {
        if (goDockerImage == null || goDockerImage.isEmpty()) {
            return DEFAULT_GO_DOCKER_IMAGE;
        }
        return goDockerImage;
    }

    public String getGoVersion() {
        if (goVersion == null || goVersion.isEmpty()) {
            return DEFAULT_GO_VERSION;
        }
        return goVersion;
    }
}
