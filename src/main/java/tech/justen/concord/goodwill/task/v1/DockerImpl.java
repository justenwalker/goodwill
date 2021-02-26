// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import com.walmartlabs.concord.sdk.ImmutableDockerContainerSpec;
import org.jetbrains.annotations.Nullable;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.DockerContainer;
import tech.justen.concord.goodwill.DockerService;

import java.io.IOException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

class DockerImpl implements DockerService {

    private static final Logger processLog = LoggerFactory.getLogger("processLog");

    private final com.walmartlabs.concord.sdk.DockerService dockerService;

    private final Context context;

    public DockerImpl(Context context, com.walmartlabs.concord.sdk.DockerService dockerService) {
        this.dockerService = dockerService;
        this.context = context;
    }

    @Override
    public int start(DockerContainer container, LogCallback outCallback, LogCallback errCallback) throws IOException, InterruptedException {
        return dockerService.start(context, ImmutableDockerContainerSpec.builder().from(new Spec(container)).build(), (String line) -> {
            processLog.debug("DOCKER: {}", line);
            if (outCallback != null) {
                outCallback.onLog(line);
            }
        }, (String line) -> {
            processLog.debug("DOCKER: {}", line);
            if (errCallback != null) {
                errCallback.onLog(line);
            }
        });
    }

    private class Spec implements com.walmartlabs.concord.sdk.DockerContainerSpec {
        private DockerContainer container;

        public Spec(DockerContainer container) {
            this.container = container;
        }

        @Override
        public String image() {
            return container.image;
        }

        @Nullable
        @Override
        public String name() {
            return container.name;
        }

        @Nullable
        @Override
        public String user() {
            return container.user;
        }

        @Nullable
        @Override
        public String workdir() {
            return container.workDir;
        }

        @Nullable
        @Override
        public String entryPoint() {
            return container.entryPoint;
        }

        @Nullable
        @Override
        public String cpu() {
            return container.cpu;
        }

        @Nullable
        @Override
        public String memory() {
            return container.memory;
        }

        @Nullable
        @Override
        public String stdOutFilePath() {
            return container.stdoutFilePath;
        }

        @Nullable
        @Override
        public List<String> args() {
            return container.command;
        }

        @Nullable
        @Override
        public Map<String, String> env() {
            return container.env;
        }

        @Nullable
        @Override
        public String envFile() {
            return container.envFile;
        }

        @Nullable
        @Override
        public Map<String, String> labels() {
            return container.labels;
        }

        @Nullable
        @Override
        public Options options() {
            if (container.hosts == null) {
                return Options.from(null);
            }
            Map<String, Object> map = new HashMap<>();
            map.put("hosts", container.hosts);
            return Options.from(map);
        }

        @Override
        public boolean debug() {
            return container.debug;
        }

        @Override
        public boolean forcePull() {
            return container.forcePull;
        }

        @Override
        public boolean redirectErrorStream() {
            return container.redirectStdError;
        }
    }
}
