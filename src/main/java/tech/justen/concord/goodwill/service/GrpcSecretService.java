// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.squareup.okhttp.Call;
import com.walmartlabs.concord.ApiClient;
import com.walmartlabs.concord.ApiException;
import com.walmartlabs.concord.ApiResponse;
import com.walmartlabs.concord.client.*;
import com.walmartlabs.concord.client.SecretOperationResponse;
import io.grpc.stub.StreamObserver;
import java.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.SecretService;
import tech.justen.concord.goodwill.TaskConfig;
import tech.justen.concord.goodwill.grpc.SecretServiceGrpc;
import tech.justen.concord.goodwill.grpc.SecretsProto;
import tech.justen.concord.goodwill.grpc.SecretsProto.*;

public class GrpcSecretService extends SecretServiceGrpc.SecretServiceImplBase {

  private static final Logger log = LoggerFactory.getLogger(GrpcSecretService.class);

  private final TaskConfig taskConfig;

  private final SecretService secretService;

  private final ApiClient apiClient;

  public GrpcSecretService(
      TaskConfig taskConfig, SecretService secretService, ApiClient apiClient) {
    this.taskConfig = taskConfig;
    this.secretService = secretService;
    this.apiClient = apiClient;
  }

  @Override
  public void createKeyPair(
      CreateKeyPairRequest request, StreamObserver<CreateKeyPairResponse> responseObserver) {
    try {
      SecretParams opts = request.getOptions();
      Map<String, Object> requestBody = createRequestBody(opts);
      requestBody.put("type", "key_pair");
      requestBody.put("private", request.getPrivateKey().toStringUtf8());
      requestBody.put("public", request.getPublicKey().toStringUtf8());
      ApiResponse<PublicKeyResponse> resp =
          create(opts.getOrgName(), requestBody, PublicKeyResponse.class);
      responseObserver.onNext(keyPairResponse(resp.getData()));
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: createKeyPair failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void generateKeyPair(
      SecretParams request, StreamObserver<CreateKeyPairResponse> responseObserver) {
    try {
      SecretParams opts = request;
      Map<String, Object> requestBody = createRequestBody(opts);
      requestBody.put("type", "key_pair");
      ApiResponse<PublicKeyResponse> resp =
          create(opts.getOrgName(), requestBody, PublicKeyResponse.class);
      responseObserver.onNext(keyPairResponse(resp.getData()));
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: generateKeyPair failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void createUsernamePassword(
      CreateUsernamePasswordRequest request,
      StreamObserver<CreateSecretResponse> responseObserver) {
    try {
      SecretParams opts = request.getOptions();
      Map<String, Object> requestBody = createRequestBody(opts);
      requestBody.put("type", "username_password");
      requestBody.put("username", request.getUsername());
      requestBody.put("password", request.getPassword());
      ApiResponse<SecretOperationResponse> resp =
          create(opts.getOrgName(), requestBody, SecretOperationResponse.class);
      responseObserver.onNext(secretResponse(resp.getData()));
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: createUsernamePassword failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void createSecretValue(
      CreateSecretValueRequest request, StreamObserver<CreateSecretResponse> responseObserver) {
    try {
      SecretParams opts = request.getOptions();
      Map<String, Object> requestBody = createRequestBody(opts);
      requestBody.put("type", "data");
      requestBody.put("data", request.getValue().toByteArray());
      ApiResponse<SecretOperationResponse> resp =
          create(opts.getOrgName(), requestBody, SecretOperationResponse.class);
      responseObserver.onNext(secretResponse(resp.getData()));
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: createSecretValue failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void exportKeyPairAsFiles(
      GetSecretRequest request, StreamObserver<KeyPairFiles> responseObserver) {
    try {
      Map<String, String> result =
          secretService.exportKeyAsFile(
              getOrg(request.getOrg()), request.getName(), request.getStorePassword());
      String pk = result.get(SecretService.PRIVATE_KEY);
      String pub = result.get(SecretService.PUBLIC_KEY);
      responseObserver.onNext(
          KeyPairFiles.newBuilder().setPrivateKeyFile(pk).setPublicKeyFile(pub).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: exportKeyPairAsFiles failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void deleteSecret(
      DeleteSecretRequest request, StreamObserver<SecretResponse> responseObserver) {
    try {
      SecretsApi secretsApi = new SecretsApi(apiClient);
      GenericOperationResult result =
          secretsApi.delete(getOrg(request.getOrg()), request.getName());
      responseObserver.onNext(SecretResponse.newBuilder().build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: deleteSecret failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void updateAccessLevels(
      UpdateSecretAccessRequest request, StreamObserver<SecretResponse> responseObserver) {
    try {
      TeamsApi teamsApi = new TeamsApi(apiClient);

      SecretsApi secretsApi = new SecretsApi(apiClient);
      List<ResourceAccessEntry> levels = new ArrayList<>();
      for (SecretsProto.AccessEntry entry : request.getEntriesList()) {
        ResourceAccessEntry e = new ResourceAccessEntry();
        if (!entry.getTeamID().isEmpty()) {
          e.setTeamId(UUID.fromString(entry.getTeamID()));
        } else {
          TeamEntry team = teamsApi.get(entry.getOrgName(), entry.getTeamName());
          e.setTeamId(team.getId());
        }
        switch (entry.getLevel()) {
          case READER:
            e.setLevel(ResourceAccessEntry.LevelEnum.READER);
          case WRITER:
            e.setLevel(ResourceAccessEntry.LevelEnum.WRITER);
          case OWNER:
            e.setLevel(ResourceAccessEntry.LevelEnum.OWNER);
        }
        levels.add(e);
      }
      GenericOperationResult result =
          secretsApi.updateAccessLevel_0(getOrg(request.getOrgName()), request.getName(), levels);
      responseObserver.onNext(SecretResponse.newBuilder().build());
      responseObserver.onCompleted();
    } catch (ApiException e) {
      log.info("ApiException: {} {}", e.getCode(), e.getResponseBody());
      log.error("GrpcSecretService: updateAccessLevels failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e, "updateAccessLevels()"));
    }
  }

  @Override
  public void listAccessLevels(
      SecretRef request, StreamObserver<ListAccessEntryResponse> responseObserver) {
    try {
      SecretsApi secretsApi = new SecretsApi(apiClient);
      List<ResourceAccessEntry> levels =
          secretsApi.getAccessLevel(request.getOrgName(), request.getName());
      ListAccessEntryResponse.Builder resp = ListAccessEntryResponse.newBuilder();
      for (ResourceAccessEntry entry : levels) {
        AccessEntry.Builder entryBuilder = AccessEntry.newBuilder();
        if (isSet(entry.getOrgName())) {
          entryBuilder.setOrgName(entry.getOrgName());
        }
        if (isSet(entry.getTeamName())) {
          entryBuilder.setTeamName(entry.getTeamName());
        }
        if (entry.getTeamId() != null) {
          entryBuilder.setTeamID(entry.getTeamId().toString());
        }
        switch (entry.getLevel()) {
          case READER:
            entryBuilder.setLevel(AccessEntry.AccessLevel.READER);
            break;
          case WRITER:
            entryBuilder.setLevel(AccessEntry.AccessLevel.WRITER);
            break;
          case OWNER:
            entryBuilder.setLevel(AccessEntry.AccessLevel.OWNER);
            break;
        }
        resp.addAccess(entryBuilder.build());
      }
      responseObserver.onNext(resp.build());
      responseObserver.onCompleted();
    } catch (ApiException e) {
      log.info("ApiException: {} {}", e.getCode(), e.getResponseBody());
      log.error("GrpcSecretService: listAccessLevels failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e, "listAccessLevels()"));
    }
  }

  @Override
  public void getUsernamePassword(
      GetSecretRequest request, StreamObserver<UsernamePassword> responseObserver) {
    try {
      Map<String, String> result =
          secretService.exportCredentials(
              getOrg(request.getOrg()), request.getName(), request.getStorePassword());
      String user = result.get(SecretService.USERNAME);
      String pass = result.get(SecretService.PASSWORD);
      responseObserver.onNext(
          UsernamePassword.newBuilder().setUsername(user).setPassword(pass).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcService: getUsernamePassword failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void exportAsFile(GetSecretRequest request, StreamObserver<SecretFile> responseObserver) {
    try {
      String filePath =
          secretService.exportAsFile(
              getOrg(request.getOrg()), request.getName(), request.getStorePassword());
      responseObserver.onNext(SecretFile.newBuilder().setFile(filePath).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: exportAsFile error", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void exportAsString(
      GetSecretRequest request, StreamObserver<SecretString> responseObserver) {
    try {
      String secret =
          secretService.exportAsString(
              getOrg(request.getOrg()), request.getName(), request.getStorePassword());
      responseObserver.onNext(SecretString.newBuilder().setStr(secret).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: exportAsString failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void decryptString(SecretString request, StreamObserver<SecretString> responseObserver) {
    try {
      String decrypted = secretService.decryptString(request.getStr());
      responseObserver.onNext(SecretString.newBuilder().setStr(decrypted).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcSecretService: decryptString failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  @Override
  public void encryptString(
      EncryptStringRequest request, StreamObserver<SecretString> responseObserver) {
    try {
      String encrypted =
          secretService.encryptString(
              getOrg(request.getOrg()), request.getProject(), request.getValue());
      responseObserver.onNext(SecretString.newBuilder().setStr(encrypted).build());
      responseObserver.onCompleted();
    } catch (Exception e) {
      log.error("GrpcService: encryptString failed", e);
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  private String getOrg(String org) {
    if (org == null || org.isEmpty()) {
      org = taskConfig.orgName();
    }
    return org;
  }

  private CreateSecretResponse secretResponse(SecretOperationResponse data) {
    CreateSecretResponse.Builder builder = CreateSecretResponse.newBuilder();
    builder.setId(data.getId().toString());
    if (data.getPassword() != null) {
      builder.setStorePassword(data.getPassword());
    }
    return builder.build();
  }

  private CreateKeyPairResponse keyPairResponse(PublicKeyResponse data) {
    CreateKeyPairResponse.Builder builder = CreateKeyPairResponse.newBuilder();
    builder.setId(data.getId().toString());
    builder.setPublicKey(data.getPublicKey());
    if (data.getPassword() != null) {
      builder.setStorePassword(data.getPassword());
    }
    return builder.build();
  }

  private <T> ApiResponse<T> create(
      String org, Map<String, Object> requestBody, Class<T> responseType) throws ApiException {
    Map<String, String> headerParams = new HashMap<>();
    headerParams.put("Content-Type", "multipart/form-data");
    if (org.isEmpty()) {
      org = taskConfig.orgName();
    }
    String url = String.format("/api/v1/org/%s/secret", org);
    String[] authNames = new String[] {"session_key", "api_key"};
    Call apiCall =
        apiClient.buildCall(
            url, "POST", null, null, null, headerParams, requestBody, authNames, null);
    return apiClient.execute(apiCall, responseType);
  }

  private Map<String, Object> createRequestBody(SecretParams opts) {
    Map<String, Object> requestBody = new HashMap<>();
    requestBody.put("name", opts.getName());
    String project = opts.getProject();
    if (!project.isEmpty()) {
      requestBody.put("project", project);
    }
    switch (opts.getVisibility()) {
      case PUBLIC:
        requestBody.put("visibility", "PUBLIC");
      case PRIVATE:
        requestBody.put("visibility", "PRIVATE");
    }
    if (opts.getGeneratePassword()) {
      requestBody.put("generatePassword", "true");
    }
    String storePassword = opts.getStorePassword();
    if (!storePassword.isEmpty()) {
      requestBody.put("storePassword", storePassword);
    }
    return requestBody;
  }

  private static boolean isSet(String str) {
    return str != null && !str.isEmpty();
  }
}
