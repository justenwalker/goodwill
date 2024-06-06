// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package tech.justen.concord.goodwill.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.json.JsonWriteFeature;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.json.JsonMapper;
import com.google.protobuf.Any;
import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Message;
import com.walmartlabs.concord.client2.ApiException;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import java.io.PrintWriter;
import java.io.StringWriter;
import java.util.*;
import tech.justen.concord.goodwill.grpc.ContextProto;
import tech.justen.concord.goodwill.grpc.ContextProto.*;

@SuppressWarnings("unchecked")
public class GrpcUtils {

  private interface Handler {
    Any getAny(Object val);
  }

  private static final ObjectMapper objectMapper =
      JsonMapper.builder().configure(JsonWriteFeature.WRITE_NUMBERS_AS_STRINGS, true).build();

  private static final Map<Class<?>, Handler> valueMapper = new HashMap<>();

  static {
    valueMapper.put(String.class, (Object obj) -> getAny((String) obj));
    valueMapper.put(Boolean.class, (Object obj) -> getAny((Boolean) obj));
    valueMapper.put(Float.class, (Object obj) -> getAny((Float) obj));
    valueMapper.put(Double.class, (Object obj) -> getAny((Double) obj));
    valueMapper.put(Byte.class, (Object obj) -> getAny((Byte) obj));
    valueMapper.put(Short.class, (Object obj) -> getAny((Short) obj));
    valueMapper.put(Integer.class, (Object obj) -> getAny((Integer) obj));
    valueMapper.put(Long.class, (Object obj) -> getAny((Long) obj));
    valueMapper.put(Date.class, (Object obj) -> getAny((Date) obj));
    valueMapper.put(Map.class, (Object obj) -> getAny((Map<String, Object>) obj));
    valueMapper.put(List.class, (Object obj) -> getAny((List<Object>) obj));
    valueMapper.put(UUID.class, (Object obj) -> getAny(((UUID) obj).toString()));
  }

  public static Object fromValue(Value val) throws InvalidProtocolBufferException {
    return fromAny(val.getValue());
  }

  public static Object fromAny(Any v) throws InvalidProtocolBufferException {
    if (v.is(NullValue.class)) {
      return null;
    }
    if (v.is(TimeValue.class)) {
      return v.unpack(TimeValue.class).getValue();
    }
    if (v.is(StringValue.class)) {
      return v.unpack(StringValue.class).getValue();
    }
    if (v.is(BoolValue.class)) {
      return v.unpack(BoolValue.class).getValue();
    }
    if (v.is(IntValue.class)) {
      return v.unpack(IntValue.class).getValue();
    }
    if (v.is(DoubleValue.class)) {
      return v.unpack(DoubleValue.class).getValue();
    }
    if (v.is(MapValue.class)) {
      Map<String, Value> map = v.unpack(MapValue.class).getValueMap();
      Map<String, Object> result = new HashMap<>();
      for (Map.Entry<String, ContextProto.Value> entry : map.entrySet()) {
        result.put(entry.getKey(), fromAny(entry.getValue().getValue()));
      }
      return result;
    }
    if (v.is(ListValue.class)) {
      List<Value> list = v.unpack(ListValue.class).getValueList();
      List<Object> result = new ArrayList<>();
      for (Value val : list) {
        result.add(fromAny(val.getValue()));
      }
      return result;
    }
    if (v.is(JSONValue.class)) {
      JSONValue val = v.unpack(JSONValue.class);
      try {
        Class<?> clazz = Class.forName(val.getClass_());
        return objectMapper.readValue(val.getJson().toStringUtf8(), clazz);
      } catch (ClassNotFoundException e) {
        throw new GrpcTypeException(val.getClass_());
      } catch (JsonMappingException e) {
        throw new GrpcValueException("JSON Mapping failed", e);
      } catch (JsonProcessingException e) {
        throw new GrpcValueException("JSON Parsing failed", e);
      }
    }
    throw new GrpcTypeException(v);
  }

  private static Any any(Message msg) {
    return Any.pack(msg);
  }

  public static Value valueOf(Object obj) {
    if (obj == null) {
      return Value.newBuilder().setValue(any(NullValue.newBuilder().build())).build();
    }
    Class<?> objClass = obj.getClass();
    for (Map.Entry<Class<?>, Handler> ent : valueMapper.entrySet()) {
      Class<?> c = ent.getKey();
      if (c.isAssignableFrom(objClass)) {
        Handler h = ent.getValue();
        return Value.newBuilder().setValue(h.getAny(obj)).build();
      }
    }
    try {
      byte[] json = objectMapper.writeValueAsBytes(obj);
      Value.newBuilder()
          .setValue(
              any(
                  JSONValue.newBuilder()
                      .setClass_(objClass.getTypeName())
                      .setJson(ByteString.copyFrom(json))
                      .build()))
          .build();
    } catch (JsonProcessingException e) {
      throw new GrpcTypeException(objClass);
    }
    throw new GrpcTypeException(objClass);
  }

  private static Any getAny(Date value) {
    return any(TimeValue.newBuilder().setValue(value.getTime()).build());
  }

  private static Any getAny(String value) {
    return any(StringValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Boolean value) {
    return any(BoolValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Float value) {
    return any(DoubleValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Double value) {
    return any(DoubleValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Byte value) {
    return any(IntValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Short value) {
    return any(IntValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Integer value) {
    return any(IntValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Long value) {
    return any(IntValue.newBuilder().setValue(value).build());
  }

  private static Any getAny(Map<String, Object> value) {
    MapValue.Builder map = MapValue.newBuilder();
    for (Map.Entry<String, Object> entry : value.entrySet()) {
      map.putValue(entry.getKey(), valueOf(entry.getValue()));
    }
    return any(map.build());
  }

  private static Any getAny(List<Object> value) {
    ListValue.Builder list = ListValue.newBuilder();
    for (Object obj : value) {
      list.addValue(valueOf(obj));
    }
    return any(list.build());
  }

  public static StatusRuntimeException toStatusException(Exception ex) {
    return Status.INTERNAL
        .withDescription(exceptionDescription(ex))
        .withCause(ex)
        .asRuntimeException();
  }

  public static StatusRuntimeException toStatusException(ApiException ex, String desc) {
    Status status = Status.INTERNAL;
    switch (ex.getCode()) {
      case 400:
        status = Status.INVALID_ARGUMENT;
        break;
      case 404:
        status = Status.NOT_FOUND;
        break;
    }
    if (!ex.getResponseBody().isEmpty()) {
      if (ex.getResponseBody().contains("already exists")) {
        status = Status.ALREADY_EXISTS;
      }
      if (ex.getResponseBody().contains("not found")) {
        status = Status.NOT_FOUND;
      }
    }
    if (desc != null && !desc.isEmpty()) {
      status = status.augmentDescription(desc);
    }
    return status.withDescription(ex.getMessage()).withCause(ex).asRuntimeException();
  }

  private static String exceptionDescription(Throwable ex) {
    StringWriter sw = new StringWriter();
    sw.append(ex.getClass().getName());
    if (ex.getMessage() != null) {
      sw.append(": ");
      sw.append(ex.getMessage());
      sw.append('\n');
    }
    ex.printStackTrace(new PrintWriter(sw));
    return sw.toString();
  }
}
