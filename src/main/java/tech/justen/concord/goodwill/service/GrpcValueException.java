// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

public class GrpcValueException extends UnsupportedOperationException {
    public GrpcValueException(String message, Throwable cause) {
        super(message, cause);
    }
}
