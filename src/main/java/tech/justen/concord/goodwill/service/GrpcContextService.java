// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import io.grpc.stub.StreamObserver;
import java.util.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.justen.concord.goodwill.ContextService;
import tech.justen.concord.goodwill.grpc.ContextProto.*;
import tech.justen.concord.goodwill.grpc.ContextServiceGrpc;

public class GrpcContextService extends ContextServiceGrpc.ContextServiceImplBase {
  private static final Logger log = LoggerFactory.getLogger(GrpcContextService.class);

  private final ContextService executionService;

  private final Map<String, Object> taskResultMap;

  public GrpcContextService(ContextService executionService, Map<String, Object> result) {
    this.executionService = executionService;
    this.taskResultMap = result;
  }

  @Override
  public void getVariable(VariableName request, StreamObserver<Value> responseObserver) {
    try {
      String name = request.getName();
      Object obj = this.executionService.getVariable(name);
      respondWithValue(responseObserver, obj);
    } catch (Exception ex) {
      log.error("GrpcExecutionService: getVariable failed", ex);
      responseObserver.onError(GrpcUtils.toStatusException(ex));
    }
  }

  @Override
  public void setVariable(Variable request, StreamObserver<SetVariableResult> responseObserver) {
    String name = request.getName();
    try {
      Object obj = GrpcUtils.fromAny(request.getValue().getValue());
      this.executionService.setVariable(name, obj);
      responseObserver.onNext(SetVariableResult.newBuilder().build());
      responseObserver.onCompleted();
    } catch (Exception ex) {
      log.error("GrpcExecutionService: setVariable failed", ex);
      responseObserver.onError(GrpcUtils.toStatusException(ex));
    }
  }

  @Override
  public void setVariables(Variables request, StreamObserver<SetVariableResult> responseObserver) {
    for (Variable var : request.getParametersList()) {
      String name = var.getName();
      try {
        Object obj = GrpcUtils.fromAny(var.getValue().getValue());
        this.executionService.setVariable(name, obj);
      } catch (Exception ex) {
        log.error(String.format("GrpcExecutionService: setVariables %s failed", name), ex);
        responseObserver.onError(GrpcUtils.toStatusException(ex));
        return;
      }
    }
    responseObserver.onNext(SetVariableResult.newBuilder().build());
    responseObserver.onCompleted();
  }

  @Override
  public void setTaskResult(Variables request, StreamObserver<SetVariableResult> responseObserver) {
    for (Variable var : request.getParametersList()) {
      String name = var.getName();
      try {
        Object obj = GrpcUtils.fromAny(var.getValue().getValue());
        this.taskResultMap.put(name, obj);
      } catch (Exception ex) {
        log.error(String.format("GrpcExecutionService: setTaskResult %s failed", name), ex);
        responseObserver.onError(GrpcUtils.toStatusException(ex));
        return;
      }
    }
    responseObserver.onNext(SetVariableResult.newBuilder().build());
    responseObserver.onCompleted();
  }

  @Override
  public void getVariableNames(
      GetVariableNameParams request, StreamObserver<VariableNameList> responseObserver) {
    try {
      Set<String> names = this.executionService.getVariableNames();
      VariableNameList ns = VariableNameList.newBuilder().addAllName(names).build();
      responseObserver.onNext(ns);
      responseObserver.onCompleted();
    } catch (Exception ex) {
      log.error("GrpcExecutionService: getVariableNames failed", ex);
      responseObserver.onError(GrpcUtils.toStatusException(ex));
    }
  }

  @Override
  public void getVariables(GetVariablesRequest request, StreamObserver<MapValue> responseObserver) {
    try {
      Set<String> names = this.executionService.getVariableNames();
      MapValue.Builder vb = MapValue.newBuilder();
      for (String name : names) {
        Object obj = this.executionService.getVariable(name);
        if (obj != null) {
          vb.putValue(name, GrpcUtils.valueOf(obj));
        }
      }
      responseObserver.onNext(vb.build());
      responseObserver.onCompleted();
    } catch (Exception ex) {
      log.error("GrpcExecutionService: getVariables failed", ex);
      responseObserver.onError(GrpcUtils.toStatusException(ex));
    }
  }

  @Override
  public void evaluate(EvaluateRequest request, StreamObserver<Value> responseObserver) {
    try {
      Class<?> cls = getClassForName(request.getType());
      // Set Parameters
      for (Variable var : request.getParametersList()) {
        Object v = GrpcUtils.fromValue(var.getValue());
        executionService.setVariable(var.getName(), v);
      }
      Object obj = executionService.evaluate(request.getExpression(), cls);
      // Unset variables
      for (Variable var : request.getParametersList()) {
        executionService.setVariable(var.getName(), null);
      }
      responseObserver.onNext(GrpcUtils.valueOf(obj));
      responseObserver.onCompleted();
    } catch (Exception ex) {
      log.error("GrpcExecutionService: evaluate failed", ex);
      responseObserver.onError(GrpcUtils.toStatusException(ex));
    }
  }

  private void respondWithValue(StreamObserver<Value> responseObserver, Object obj) {
    if (obj == null) {
      responseObserver.onCompleted();
      return;
    }
    try {
      Value val = GrpcUtils.valueOf(obj);
      responseObserver.onNext(val);
      responseObserver.onCompleted();
    } catch (Exception e) {
      responseObserver.onError(GrpcUtils.toStatusException(e));
    }
  }

  private Class<?> getClassForName(String name) throws ClassNotFoundException {
    switch (name) {
      case "time":
        return Date.class;
      case "map":
        return Map.class;
      case "list":
        return List.class;
      case "string":
        return String.class;
      case "int32":
        return Integer.class;
      case "float32":
        return Float.class;
      case "float":
      case "float64":
        return Double.class;
      case "int":
      case "int64":
        return Long.class;
      case "bool":
        return Boolean.class;
      case "json":
        return Object.class;
      case "null":
        return Object.class;
    }
    throw new ClassNotFoundException(
        String.format("cannot find java class for value type %s", name));
  }
}
