// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.walmartlabs.concord.client2.ApiClient;
import io.grpc.stub.StreamObserver;
import tech.justen.concord.goodwill.TaskConfig;
import tech.justen.concord.goodwill.grpc.ConfigProto.*;
import tech.justen.concord.goodwill.grpc.ConfigServiceGrpc;

public class GrpcConfigService extends ConfigServiceGrpc.ConfigServiceImplBase {

  private final TaskConfig taskConfig;

  private final ApiClient apiClient;

  private final String sessionToken;

  public GrpcConfigService(ApiClient apiClient, TaskConfig taskConfig, String sessionToken) {
    this.apiClient = apiClient;
    this.taskConfig = taskConfig;
    this.sessionToken = sessionToken;
  }

  @Override
  public void getConfiguration(
      ConfigurationRequest request, StreamObserver<Configuration> responseObserver) {
    responseObserver.onNext(
        Configuration.newBuilder()
            .setProcessID(taskConfig.processId())
            .setWorkingDirectory(taskConfig.workingDirectory().toString())
            .setApiConfiguration(
                APIConfiguration.newBuilder()
                    .setBaseURL(apiClient.getBaseUrl())
                    .setSessionToken(sessionToken)
                    .build())
            .setProjectInfo(
                ProjectInfo.newBuilder()
                    .setOrgID(taskConfig.orgId())
                    .setOrgName(taskConfig.orgName())
                    .setProjectID(taskConfig.projectId())
                    .setProjectName(taskConfig.projectName())
                    .build())
            .setRepositoryInfo(
                RepositoryInfo.newBuilder()
                    .setRepoID(taskConfig.repoId())
                    .setRepoName(taskConfig.repoName())
                    .setRepoURL(taskConfig.repoUrl())
                    .build())
            .build());
    responseObserver.onCompleted();
  }
}
