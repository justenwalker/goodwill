// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import io.grpc.stub.StreamObserver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.LockService;
import tech.justen.concord.goodwill.grpc.LockProto;
import tech.justen.concord.goodwill.grpc.LockServiceGrpc;

public class GrpcLockService extends LockServiceGrpc.LockServiceImplBase {
    private static final Logger log = LoggerFactory.getLogger(GrpcLockService.class);

    private final LockService lockService;

    public GrpcLockService(LockService lockService) {
        this.lockService = lockService;
    }

    @Override
    public void projectLock(LockProto.Lock request, StreamObserver<LockProto.LockResult> responseObserver) {
        try {
            lockService.projectLock(request.getName());
            responseObserver.onNext(LockProto.LockResult.newBuilder().build());
            responseObserver.onCompleted();
        } catch (Exception ex) {
            log.error("GrpcLockService: projectLock failed", ex);
            responseObserver.onError(GrpcUtils.toStatusException(ex));
        }
    }

    @Override
    public void projectUnlock(LockProto.Lock request, StreamObserver<LockProto.LockResult> responseObserver) {
        try {
            lockService.projectUnlock(request.getName());
            responseObserver.onNext(LockProto.LockResult.newBuilder().build());
            responseObserver.onCompleted();
        } catch (Exception ex) {
            log.error("GrpcLockService: projectUnlock failed", ex);
            responseObserver.onError(GrpcUtils.toStatusException(ex));
        }
    }
}
