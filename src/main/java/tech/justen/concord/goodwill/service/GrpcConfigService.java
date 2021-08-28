// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.walmartlabs.concord.client.ApiClientConfiguration;
import io.grpc.stub.StreamObserver;
import tech.justen.concord.goodwill.BuildInfo;
import tech.justen.concord.goodwill.TaskConfig;
import tech.justen.concord.goodwill.grpc.ConfigProto.*;
import tech.justen.concord.goodwill.grpc.ConfigServiceGrpc;

public class GrpcConfigService extends ConfigServiceGrpc.ConfigServiceImplBase {

    private final TaskConfig taskConfig;

    private final ApiClientConfiguration apiClientConfig;

    public GrpcConfigService(ApiClientConfiguration apiClientConfig, TaskConfig taskConfig) {
        this.apiClientConfig = apiClientConfig;
        this.taskConfig = taskConfig;
    }

    @Override
    public void getConfiguration(ConfigurationRequest request, StreamObserver<Configuration> responseObserver) {
        responseObserver.onNext(Configuration.newBuilder()
                .setProcessID(taskConfig.processId())
                .setWorkingDirectory(taskConfig.workingDirectory().toString())
                .setApiConfiguration(APIConfiguration.newBuilder()
                        .setBaseURL(apiClientConfig.baseUrl())
                        .setSessionToken(apiClientConfig.sessionToken())
                        .build())
                .setProjectInfo(ProjectInfo.newBuilder()
                        .setOrgID(taskConfig.orgId())
                        .setOrgName(taskConfig.orgName())
                        .setProjectID(taskConfig.projectId())
                        .setProjectName(taskConfig.projectName())
                        .build())
                .setRepositoryInfo(RepositoryInfo.newBuilder()
                        .setRepoID(taskConfig.repoId())
                        .setRepoName(taskConfig.repoName())
                        .setRepoURL(taskConfig.repoUrl())
                        .build())
                .build());
        responseObserver.onCompleted();
    }
}
