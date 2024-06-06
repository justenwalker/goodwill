// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import com.walmartlabs.concord.sdk.ContextUtils;
import com.walmartlabs.concord.sdk.ProjectInfo;
import java.util.Map;
import tech.justen.concord.goodwill.SecretService;

public class SecretServiceImpl implements SecretService {

  private final com.walmartlabs.concord.sdk.SecretService secretService;

  private final Context ctx;

  public SecretServiceImpl(Context ctx, com.walmartlabs.concord.sdk.SecretService secretService) {
    this.ctx = ctx;
    this.secretService = secretService;
  }

  @Override
  public String exportAsString(String orgName, String name, String password) throws Exception {
    return secretService.exportAsString(ctx, getInstanceId(), getOrgName(orgName), name, password);
  }

  @Override
  public Map<String, String> exportKeyAsFile(String orgName, String name, String password)
      throws Exception {
    return secretService.exportKeyAsFile(
        ctx, getInstanceId(), getWorkDir(), getOrgName(orgName), name, password);
  }

  @Override
  public Map<String, String> exportCredentials(String orgName, String name, String password)
      throws Exception {
    return secretService.exportCredentials(
        ctx, getInstanceId(), getWorkDir(), getOrgName(orgName), name, password);
  }

  @Override
  public String exportAsFile(String orgName, String name, String password) throws Exception {
    return secretService.exportAsFile(
        ctx, getInstanceId(), getWorkDir(), getOrgName(orgName), name, password);
  }

  @Override
  public String decryptString(String s) throws Exception {
    String instanceId = ContextUtils.getTxId(ctx).toString();
    return secretService.decryptString(ctx, instanceId, s);
  }

  @Override
  public String encryptString(String orgName, String projectName, String value) throws Exception {
    return secretService.encryptString(
        ctx, getInstanceId(), getOrgName(orgName), getProjectName(projectName), value);
  }

  private String getInstanceId() {
    return ContextUtils.getTxId(ctx).toString();
  }

  private String getWorkDir() {
    return ContextUtils.getWorkDir(ctx).toString();
  }

  private String getOrgName(String orgName) {
    if (orgName != null && !orgName.isEmpty()) {
      return orgName;
    }
    ProjectInfo pinfo = ContextUtils.getProjectInfo(ctx);
    if (pinfo == null) {
      return "Default";
    }
    return pinfo.orgName();
  }

  private String getProjectName(String projectName) {
    if (projectName != null && !projectName.isEmpty()) {
      return projectName;
    }
    ProjectInfo pinfo = ContextUtils.getProjectInfo(ctx);
    if (pinfo == null) {
      return "";
    }
    return pinfo.name();
  }
}
