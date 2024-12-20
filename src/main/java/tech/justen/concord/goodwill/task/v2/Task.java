// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v2;

import static tech.justen.concord.goodwill.task.TaskParams.*;

import com.walmartlabs.concord.client2.ApiClient;
import com.walmartlabs.concord.runtime.v2.sdk.Context;
import com.walmartlabs.concord.runtime.v2.sdk.DependencyManager;
import com.walmartlabs.concord.runtime.v2.sdk.TaskResult;
import com.walmartlabs.concord.runtime.v2.sdk.Variables;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import javax.inject.Inject;
import javax.inject.Named;
import tech.justen.concord.goodwill.task.TaskCommon;
import tech.justen.concord.goodwill.task.TaskParams;

@Named("goodwill")
public class Task implements com.walmartlabs.concord.runtime.v2.sdk.Task {

  private final DependencyManager dependencyManager;

  private final ExecutorService executor;

  private final Context context;

  private final ApiClient apiClient;

  private final Variables defaults;

  @Inject
  public Task(Context context, DependencyManager dependencyManager, ApiClient apiClient) {
    this.defaults = context.defaultVariables();
    this.dependencyManager = dependencyManager;
    this.executor = Executors.newCachedThreadPool();
    this.context = context;
    this.apiClient = apiClient;
  }

  @Override
  public TaskResult execute(Variables input) throws Exception {
    String sessionToken = context.processConfiguration().processInfo().sessionToken();
    TaskCommon common =
        new TaskCommon(
            new TaskConfigImpl(context),
            makeParams(input),
            new DependencyManagerImpl(dependencyManager),
            new ContextServiceImpl(context, input),
            new DockerImpl(context.dockerService()),
            new SecretServiceImpl(context),
            new LockServiceImpl(context.lockService()),
            executor,
            apiClient,
            sessionToken);
    Map<String, Object> result = common.execute();
    return TaskResult.of(true).values(result);
  }

  private String defaultString(Variables input, String key) {
    if (!input.has(key)) {
      return defaults.getString(key);
    }
    return input.getString(key);
  }

  private boolean defaultBool(Variables input, String key, boolean defaultValue) {
    if (!input.has(key)) {
      return defaults.getBoolean(key, defaultValue);
    }
    return input.getBoolean(key, defaultValue);
  }

  private TaskParams makeParams(Variables input) {
    TaskParams params = new TaskParams();
    params.goArch = defaultString(input, GOARCH_KEY);
    params.goOS = defaultString(input, GOOS_KEY);
    params.goDockerImage = defaultString(input, GO_DOCKER_IMAGE_KEY);
    params.goVersion = defaultString(input, GO_VERSION_KEY);
    params.installGo = defaultBool(input, INSTALL_GO_KEY, false);
    params.useDockerImage = defaultBool(input, USE_DOCKER_IMAGE_KEY, true);
    params.debug = defaultBool(input, DEBUG_KEY, false);
    params.tasksBinary = defaultString(input, BINARY_KEY);
    params.buildDir = defaultString(input, BUILD_DIR_KEY);
    params.taskName = defaultString(input, TASK_NAME_KEY);
    params.goProxy = defaultString(input, GOPROXY_KEY);
    params.goNoProxy = defaultString(input, GONOPROXY_KEY);
    params.goPrivate = defaultString(input, GOPRIVATE_KEY);
    params.goSumDB = defaultString(input, GOSUMDB_KEY);
    params.goNoSumDB = defaultString(input, GONOSUMDB_KEY);
    return params;
  }
}
