// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill;

public interface LockService {
    void projectLock(String name) throws Exception;

    void projectUnlock(String name) throws Exception;
}
