// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import com.walmartlabs.concord.sdk.ContextUtils;
import com.walmartlabs.concord.sdk.ProjectInfo;
import com.walmartlabs.concord.sdk.RepositoryInfo;
import java.nio.file.Path;
import java.util.UUID;
import tech.justen.concord.goodwill.TaskConfig;

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
    return notNull(txId);
  }

  @Override
  public String orgId() {
    return notNull(orgId);
  }

  @Override
  public String projectName() {
    return notNull(projectName);
  }

  @Override
  public String projectId() {
    return notNull(projectId);
  }

  @Override
  public String repoName() {
    return notNull(repoName);
  }

  @Override
  public String repoId() {
    return notNull(repoId);
  }

  @Override
  public String repoUrl() {
    return notNull(repoURL);
  }

  @Override
  public String orgName() {
    return notNull(orgName);
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

  private static String notNull(String str) {
    if (str == null) {
      return "";
    }
    return str;
  }
}
