// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v2;

import com.walmartlabs.concord.runtime.v2.sdk.Context;
import com.walmartlabs.concord.runtime.v2.sdk.ProjectInfo;
import tech.justen.concord.goodwill.SecretService;

import java.util.HashMap;
import java.util.Map;

public class SecretServiceImpl implements SecretService {

    private final com.walmartlabs.concord.runtime.v2.sdk.SecretService secretService;

    private final Context context;

    public SecretServiceImpl(Context context) {
        this.context = context;
        this.secretService = context.secretService();
    }

    private String getOrgName(String orgName) {
        if (orgName != null && !orgName.isEmpty()) {
            return orgName;
        }
        ProjectInfo info = context.processConfiguration().projectInfo();
        return info.orgName();
    }

    private String getProjectName(String projectName) {
        if (projectName != null && !projectName.isEmpty()) {
            return projectName;
        }
        ProjectInfo info = context.processConfiguration().projectInfo();
        return info.projectName();
    }

    @Override
    public String exportAsString(String orgName, String name, String password) throws Exception {
        return secretService.exportAsString(getOrgName(orgName), name, password);
    }

    @Override
    public Map<String, String> exportKeyAsFile(String orgName, String name, String password) throws Exception {
        com.walmartlabs.concord.runtime.v2.sdk.SecretService.KeyPair kf = secretService.exportKeyAsFile(getOrgName(orgName), name, password);
        Map<String, String> m = new HashMap<>();
        m.put(PUBLIC_KEY, kf.publicKey().toString());
        m.put(PRIVATE_KEY, kf.publicKey().toString());
        return m;
    }

    @Override
    public Map<String, String> exportCredentials(String orgName, String name, String password) throws Exception {
        com.walmartlabs.concord.runtime.v2.sdk.SecretService.UsernamePassword up = secretService.exportCredentials(getOrgName(orgName), name, password);
        Map<String, String> m = new HashMap<>();
        m.put(USERNAME, up.username());
        m.put(PASSWORD, up.password());
        return m;
    }

    @Override
    public String exportAsFile(String orgName, String name, String password) throws Exception {
        return secretService.exportAsFile(getOrgName(orgName), name, password).toString();
    }

    @Override
    public String decryptString(String s) throws Exception {
        return secretService.decryptString(s);
    }

    @Override
    public String encryptString(String orgName, String projectName, String value) throws Exception {
        return secretService.encryptString(getOrgName(orgName), getProjectName(projectName), value);
    }
}
