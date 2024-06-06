// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.fasterxml.jackson.core.type.TypeReference;
import com.walmartlabs.concord.client2.*;
import com.walmartlabs.concord.client2.impl.HttpEntity;
import com.walmartlabs.concord.client2.impl.MultipartRequestBodyHandler;
import com.walmartlabs.concord.client2.impl.ResponseBodyHandler;
import com.walmartlabs.concord.sdk.Constants;
import io.grpc.stub.StreamObserver;
import java.io.IOException;
import java.io.InputStream;
import java.net.URI;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.SecretService;
import tech.justen.concord.goodwill.TaskConfig;
import tech.justen.concord.goodwill.grpc.SecretServiceGrpc;
import tech.justen.concord.goodwill.grpc.SecretsProto.*;

public class GrpcSecretService extends SecretServiceGrpc.SecretServiceImplBase {

  private static final Logger log = LoggerFactory.getLogger(GrpcSecretService.class);

  private final TaskConfig taskConfig;

  private final SecretService secretService;

  private final ApiClient apiClient;

  private final SecretClient secretClient;

  public GrpcSecretService(
      TaskConfig taskConfig, SecretService secretService, ApiClient apiClient) {
    this.taskConfig = taskConfig;
    this.secretService = secretService;
    this.apiClient = apiClient;
    this.secretClient = new SecretClient(apiClient);
  }

  @Override
  public void createKeyPair(
      CreateKeyPairRequest request, StreamObserver<CreateKeyPairResponse> responseObserver) {
    try {
      CreateKeyPairResponse response = createKeyPairRaw(request);
      responseObserver.onNext(response);
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
      CreateKeyPairResponse response =
          createKeyPairRaw(CreateKeyPairRequest.newBuilder().setOptions(request).build());
      responseObserver.onNext(response);
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
      ImmutableCreateSecretRequest.Builder requestBuilder =
          createSecretRequestBuilder(request.getOptions());
      requestBuilder.usernamePassword(
          CreateSecretRequest.UsernamePassword.of(request.getUsername(), request.getPassword()));
      SecretOperationResponse response = secretClient.createSecret(requestBuilder.build());
      responseObserver.onNext(secretResponse(response));
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
      ImmutableCreateSecretRequest.Builder requestBuilder =
          createSecretRequestBuilder(request.getOptions());
      requestBuilder.data(request.getValue().toByteArray());
      SecretOperationResponse response = secretClient.createSecret(requestBuilder.build());
      responseObserver.onNext(secretResponse(response));
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
      for (AccessEntry entry : request.getEntriesList()) {
        ResourceAccessEntry e = new ResourceAccessEntry();
        if (!entry.getTeamID().isEmpty()) {
          e.setTeamId(UUID.fromString(entry.getTeamID()));
        } else {
          TeamEntry team = teamsApi.getTeam(entry.getOrgName(), entry.getTeamName());
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
      secretsApi.updateSecretAccessLevelBulk(
          getOrg(request.getOrgName()), request.getName(), levels);
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
          secretsApi.getSecretAccessLevel(request.getOrgName(), request.getName());
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

  private static final TypeReference<PublicKeyResponse> PUBLIC_KEY_RESPONSE =
      new TypeReference<PublicKeyResponse>() {};

  /**
   * Creates a keypair secret with raw bytes. The builder only allows KeyPairs to be provided as
   * files on the filesystem; so in order to avoid temporary files, this needs to do a manual api
   * call using the bytes from CreateKeyPairRequest.
   *
   * @param request is a CreateKeyPairRequest to create the KeyPair secret.
   * @return CreateKeyPairResponse from the API
   * @throws ApiException when it encounters an API Error
   */
  private CreateKeyPairResponse createKeyPairRaw(CreateKeyPairRequest request) throws ApiException {
    SecretParams opts = request.getOptions();
    String orgName = opts.getOrgName();
    if (orgName.isBlank()) {
      orgName = taskConfig.orgName();
    }
    Map<String, Object> params = new HashMap<>();
    params.put(Constants.Multipart.NAME, opts.getName());
    params.put(Constants.Multipart.GENERATE_PASSWORD, opts.getGeneratePassword());
    String storePassword = opts.getStorePassword();
    if (!storePassword.isEmpty()) {
      params.put(Constants.Multipart.STORE_PASSWORD, storePassword);
    }
    switch (opts.getVisibility()) {
      case PRIVATE:
        params.put(Constants.Multipart.VISIBILITY, SecretEntryV2.VisibilityEnum.PRIVATE.getValue());
        break;
      case PUBLIC:
        params.put(Constants.Multipart.VISIBILITY, SecretEntryV2.VisibilityEnum.PUBLIC.getValue());
        break;
    }
    String project = opts.getProject();
    if (!project.isEmpty()) {
      params.put(Constants.Multipart.PROJECT_NAMES, project);
    }
    params.put(Constants.Multipart.TYPE, SecretEntryV2.TypeEnum.KEY_PAIR.getValue());
    if (!request.getPrivateKey().isEmpty() && !request.getPublicKey().isEmpty()) {
      params.put(Constants.Multipart.PUBLIC, request.getPublicKey().toByteArray());
      params.put(Constants.Multipart.PRIVATE, request.getPrivateKey().toByteArray());
    }
    HttpRequest.Builder requestBuilder = apiClient.requestBuilder();
    String createSecretPath =
        "/api/v1/org/{orgName}/secret".replace("{orgName}", ApiClient.urlEncode(orgName));
    requestBuilder.uri(URI.create(apiClient.getBaseUri() + createSecretPath));
    requestBuilder.header(
        "Accept", "application/json,application/vnd.siesta-validation-errors-v1+json");
    HttpEntity entity = MultipartRequestBodyHandler.handle(apiClient.getObjectMapper(), params);
    requestBuilder
        .header("Content-Type", entity.contentType().toString())
        .method(
            "POST",
            HttpRequest.BodyPublishers.ofInputStream(
                () -> {
                  try {
                    return entity.getContent();
                  } catch (IOException e) {
                    throw new RuntimeException(e);
                  }
                }));
    try {
      HttpResponse<InputStream> response =
          apiClient
              .getHttpClient()
              .send(requestBuilder.build(), HttpResponse.BodyHandlers.ofInputStream());
      if (response.statusCode() / 100 != 2) {
        try (InputStream body = response.body()) {
          String bodyString = body == null ? null : new String(body.readAllBytes());
          String message =
              formatExceptionMessage("createSecret", response.statusCode(), bodyString);
          throw new ApiException(response.statusCode(), message, response.headers(), bodyString);
        }
      }
      return keyPairResponse(
          ResponseBodyHandler.handle(apiClient.getObjectMapper(), response, PUBLIC_KEY_RESPONSE));
    } catch (IOException e) {
      throw new ApiException(e);
    } catch (InterruptedException e) {
      Thread.currentThread().interrupt();
      throw new ApiException(e);
    }
  }

  private ImmutableCreateSecretRequest.Builder createSecretRequestBuilder(SecretParams opts) {
    ImmutableCreateSecretRequest.Builder requestBuilder = CreateSecretRequest.builder();
    requestBuilder.name(opts.getName());
    requestBuilder.org(opts.getOrgName());
    if (opts.getOrgName().isBlank()) {
      requestBuilder.org(taskConfig.orgName());
    }
    String project = opts.getProject();
    if (!project.isEmpty()) {
      requestBuilder.addProjectNames(project);
    }
    switch (opts.getVisibility()) {
      case PUBLIC:
        requestBuilder.visibility(SecretEntryV2.VisibilityEnum.PUBLIC);
      case PRIVATE:
        requestBuilder.visibility(SecretEntryV2.VisibilityEnum.PRIVATE);
    }
    if (opts.getGeneratePassword()) {
      requestBuilder.generatePassword(true);
    }
    String storePassword = opts.getStorePassword();
    if (!storePassword.isEmpty()) {
      requestBuilder.storePassword(storePassword);
    }
    return requestBuilder;
  }

  private static boolean isSet(String str) {
    return str != null && !str.isEmpty();
  }

  private static String formatExceptionMessage(String operationId, int statusCode, String body) {
    if (body == null || body.isEmpty()) {
      body = "[no body]";
    }
    return operationId + " call failed with: " + statusCode + " - " + body;
  }
}
