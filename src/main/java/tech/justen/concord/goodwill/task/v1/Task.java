// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import static tech.justen.concord.goodwill.task.TaskParams.*;

import com.walmartlabs.concord.client2.ApiClient;
import com.walmartlabs.concord.client2.ApiClientFactory;
import com.walmartlabs.concord.client2.ImmutableApiClientConfiguration;
import com.walmartlabs.concord.sdk.*;
import java.util.Collections;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import javax.inject.Inject;
import javax.inject.Named;
import tech.justen.concord.goodwill.task.TaskCommon;
import tech.justen.concord.goodwill.task.TaskParams;

@Named("goodwill")
public class Task implements com.walmartlabs.concord.sdk.Task {

  private final DependencyManager dependencyManager;

  private final SecretService secretService;

  private final LockService lockService;

  private final ExecutorService executor;

  private final DockerService dockerService;

  private final ApiClientFactory apiClientFactory;

  private final ApiConfiguration apiConfiguration;

  private Map<String, Object> defaults;

  @Inject
  public Task(
      DependencyManager dependencyManager,
      DockerService dockerService,
      SecretService secretService,
      LockService lockService,
      ApiClientFactory apiClientFactory,
      ApiConfiguration apiConfiguration) {
    this.dockerService = dockerService;
    this.dependencyManager = dependencyManager;
    this.secretService = secretService;
    this.lockService = lockService;
    this.executor = Executors.newCachedThreadPool();
    this.apiClientFactory = apiClientFactory;
    this.apiConfiguration = apiConfiguration;
  }

  @Override
  public void execute(Context ctx) throws java.lang.Exception {
    defaults = ContextUtils.getMap(ctx, "goodwillCfg");
    if (defaults == null) {
      defaults = Collections.emptyMap();
    }
    String sessionToken = ContextUtils.getSessionToken(ctx);
    ApiClient apiClient =
        apiClientFactory.create(
            ImmutableApiClientConfiguration.builder()
                .baseUrl(apiConfiguration.getBaseUrl())
                .sessionToken(sessionToken)
                .build());
    TaskCommon common =
        new TaskCommon(
            new TaskConfigImpl(ctx),
            makeParams(ctx),
            new DependencyManagerImpl(dependencyManager),
            new ContextServiceImpl(ctx),
            new DockerImpl(ctx, dockerService),
            new SecretServiceImpl(ctx, secretService),
            new LockServiceImpl(ctx, lockService),
            executor,
            apiClient,
            sessionToken);
    Map<String, Object> result = common.execute();
    for (Map.Entry<String, Object> e : result.entrySet()) {
      ctx.setVariable(e.getKey(), e.getValue());
    }
  }

  private String defaultString(Context ctx, String key) {
    String value = ContextUtils.getString(ctx, key);
    if (value == null) {
      return MapUtils.getString(defaults, key);
    }
    return value;
  }

  private boolean defaultBool(Context ctx, String key, boolean defaultValue) {
    return ContextUtils.getBoolean(ctx, key, MapUtils.getBoolean(defaults, key, defaultValue));
  }

  private TaskParams makeParams(Context ctx) {
    TaskParams params = new TaskParams();
    params.goArch = defaultString(ctx, GOARCH_KEY);
    params.goOS = defaultString(ctx, GOOS_KEY);
    params.goDockerImage = defaultString(ctx, GO_DOCKER_IMAGE_KEY);
    params.goVersion = defaultString(ctx, GO_VERSION_KEY);
    params.installGo = defaultBool(ctx, INSTALL_GO_KEY, true);
    params.useDockerImage = defaultBool(ctx, USE_DOCKER_IMAGE_KEY, false);
    params.debug = defaultBool(ctx, DEBUG_KEY, false);
    params.tasksBinary = defaultString(ctx, BINARY_KEY);
    params.buildDir = defaultString(ctx, BUILD_DIR_KEY);
    params.taskName = defaultString(ctx, TASK_NAME_KEY);
    params.goProxy = defaultString(ctx, GOPROXY_KEY);
    params.goNoProxy = defaultString(ctx, GONOPROXY_KEY);
    params.goPrivate = defaultString(ctx, GOPRIVATE_KEY);
    params.goSumDB = defaultString(ctx, GOSUMDB_KEY);
    params.goNoSumDB = defaultString(ctx, GONOSUMDB_KEY);
    return params;
  }
}
