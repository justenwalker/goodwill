// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v2;

import tech.justen.concord.goodwill.LockService;

public class LockServiceImpl implements LockService {
  private final com.walmartlabs.concord.runtime.v2.sdk.LockService lockService;

  public LockServiceImpl(com.walmartlabs.concord.runtime.v2.sdk.LockService lockService) {
    this.lockService = lockService;
  }

  @Override
  public void projectLock(String name) throws Exception {
    lockService.projectLock(name);
  }

  @Override
  public void projectUnlock(String name) throws Exception {
    lockService.projectUnlock(name);
  }
}
