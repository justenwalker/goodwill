// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v2;

import com.walmartlabs.concord.runtime.v2.sdk.DependencyManager;
import java.io.IOException;
import java.net.URI;
import java.nio.file.Path;

class DependencyManagerImpl implements tech.justen.concord.goodwill.task.DependencyManager {

  private final DependencyManager dependencyManager;

  public DependencyManagerImpl(DependencyManager dependencyManager) {
    this.dependencyManager = dependencyManager;
  }

  @Override
  public Path resolve(URI uri) throws IOException {
    return dependencyManager.resolve(uri);
  }
}
