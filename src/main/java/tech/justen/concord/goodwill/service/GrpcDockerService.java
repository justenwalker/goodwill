// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import io.grpc.stub.StreamObserver;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.DockerContainer;
import tech.justen.concord.goodwill.DockerService;
import tech.justen.concord.goodwill.grpc.DockerProto.*;
import tech.justen.concord.goodwill.grpc.DockerServiceGrpc;

public class GrpcDockerService extends DockerServiceGrpc.DockerServiceImplBase {

    private static final Logger log = LoggerFactory.getLogger(GrpcDockerService.class);

    private final DockerService dockerService;

    public GrpcDockerService(DockerService dockerService) {
        this.dockerService = dockerService;
    }

    @Override
    public void runContainer(DockerContainerSpec request, StreamObserver<DockerContainerResult> responseObserver) {
        try {
            int result = dockerService.start(DockerContainer.fromGrpcRequest(request),
                    (String line) -> {
                        responseObserver.onNext(DockerContainerResult.newBuilder().setStdout(line).build());
                    }, (String line) -> {
                        responseObserver.onNext(DockerContainerResult.newBuilder().setStderr(line).build());
                    });
            responseObserver.onNext(DockerContainerResult.newBuilder().setStatus(result).build());
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("GrpcService: runContainer error", e);
            responseObserver.onError(GrpcUtils.toStatusException(e));
        }
    }
}
