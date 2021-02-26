// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import com.walmartlabs.concord.sdk.ContextUtils;
import com.walmartlabs.concord.sdk.ProjectInfo;
import com.walmartlabs.concord.sdk.RepositoryInfo;
import tech.justen.concord.goodwill.TaskConfig;

import java.nio.file.Path;
import java.util.UUID;

public class TaskConfigImpl implements TaskConfig {

    private final Context ctx;

    private final String txId;

    private final Path workingDirectory;

    private final String orgName;

    private final String orgId;

    private final String projectName;

    private final String projectId;

    private final String repoId;

    private final String repoName;

    private final String repoURL;

    public TaskConfigImpl(Context ctx) {
        this.ctx = ctx;
        String txId = ContextUtils.getTxId(ctx).toString();
        Path workingDirectory = ContextUtils.getWorkDir(ctx);
        String orgId = "";
        String orgName = "Default";
        String projectName = "";
        String projectId = "";
        ProjectInfo projectInfo = ContextUtils.getProjectInfo(ctx);
        if (projectInfo != null) {
            orgId = uuidStr(projectInfo.orgId());
            orgName = projectInfo.orgName();
            projectId = uuidStr(projectInfo.id());
            projectName = projectInfo.name();
        }
        String repoName = "";
        String repoId = "";
        String repoURL = "";
        RepositoryInfo repoInfo = ContextUtils.getRepositoryInfo(ctx);
        if (repoInfo != null) {
            repoName = repoInfo.name();
            repoId = uuidStr(repoInfo.id());
            repoURL = repoInfo.url();
        }
        this.txId = txId;
        this.workingDirectory = workingDirectory;
        this.orgId = orgId;
        this.orgName = orgName;
        this.projectName = projectName;
        this.projectId = projectId;
        this.repoId = repoId;
        this.repoName = repoName;
        this.repoURL = repoURL;
    }

    @Override
    public String processId() {
        return txId;
    }

    @Override
    public String orgId() {
        return orgId;
    }

    @Override
    public String projectName() {
        return projectName;
    }

    @Override
    public String projectId() {
        return projectId;
    }

    @Override
    public String repoName() {
        return repoName;
    }

    @Override
    public String repoId() {
        return repoId;
    }

    @Override
    public String repoUrl() {
        return repoURL;
    }

    @Override
    public String orgName() {
        return orgName;
    }

    @Override
    public Path workingDirectory() {
        return workingDirectory;
    }

    private static String uuidStr(UUID uuid) {
        if (uuid == null) {
            return "";
        }
        return uuid.toString();
    }
}
