// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.google.protobuf.Any;

public class GrpcValueException extends UnsupportedOperationException {
    public GrpcValueException(Class clazz) {
        super(String.format("Unsupported Java Class: %s", clazz.getName()));
    }

    public GrpcValueException(Any any) {
        super(String.format("Unsupported Any Type: %s", any.getDescriptorForType().getFullName()));
    }
}
