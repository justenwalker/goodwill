// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.google.protobuf.Any;

public class GrpcTypeException extends UnsupportedOperationException {
    public GrpcTypeException(Class clazz) {
        this(clazz.getName());
    }

    public GrpcTypeException(Any any) {
        super(String.format("Unsupported Any Type: %s", any.getTypeUrl()));
    }

    public GrpcTypeException(String className) {
        super(String.format("Unsupported Java Class: %s", className));
    }
}
