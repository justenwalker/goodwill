// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

import tech.justen.concord.goodwill.grpc.DockerProto.*;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class DockerContainer {

    public static final String DEFAULT_WORK_DIR = "/workspace";

    public static DockerContainer fromGrpcRequest(DockerContainerSpec spec) {
        DockerContainer container = new DockerContainer();
        container.workDir = DEFAULT_WORK_DIR;
        if (!spec.getImage().isEmpty()) {
            container.image = spec.getImage();
        }
        if (!spec.getName().isEmpty()) {
            container.name = spec.getName();
        }
        if (!spec.getUser().isEmpty()) {
            container.user = spec.getUser();
        }
        if (!spec.getWorkDir().isEmpty()) {
            container.workDir = spec.getWorkDir();
        }
        if (!spec.getEntryPoint().isEmpty()) {
            container.entryPoint = spec.getEntryPoint();
        }
        if (spec.getCommandCount() > 0) {
            container.command = new ArrayList<>(spec.getCommandList());
        }
        if (!spec.getCpu().isEmpty()) {
            container.cpu = spec.getCpu();
        }
        if (!spec.getMemory().isEmpty()) {
            container.memory = spec.getMemory();
        }
        if (spec.getEnvCount() > 0) {
            container.env = new HashMap<>(spec.getEnvMap());
        }
        if (!spec.getEnvFile().isEmpty()) {
            container.envFile = spec.getEnvFile();
        }
        if (spec.getLabelsCount() > 0) {
            container.labels = new HashMap<>(spec.getLabelsMap());
        }
        container.forcePull = spec.getForcePull();
        if (!spec.getStdoutFilePath().isEmpty()) {
            container.stdoutFilePath = spec.getStdoutFilePath();
        }
        container.redirectStdError = spec.getRedirectStdError();
        return container;
    }

    public String image;

    public String name;

    public String user;

    public String workDir;

    public String entryPoint;

    public List<String> command;

    public String cpu;

    public String memory;

    public Map<String, String> env;

    public String envFile;

    public Map<String, String> labels;

    public boolean forcePull;

    public List<String> hosts;

    public String stdoutFilePath;

    public boolean redirectStdError;

    public boolean debug;
}
