// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.task.v1;

import com.walmartlabs.concord.sdk.Context;
import tech.justen.concord.goodwill.LockService;

public class LockServiceImpl implements LockService {
    private final Context ctx;

    private final com.walmartlabs.concord.sdk.LockService lockService;

    public LockServiceImpl(Context ctx, com.walmartlabs.concord.sdk.LockService lockService) {
        this.ctx = ctx;
        this.lockService = lockService;
    }

    @Override
    public void projectLock(String name) throws Exception {
        lockService.projectLock(ctx, name);
    }

    @Override
    public void projectUnlock(String name) throws Exception {
        lockService.projectUnlock(ctx, name);
    }
}
