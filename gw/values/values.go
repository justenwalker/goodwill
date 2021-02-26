// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"encoding/json"
	"fmt"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
)

type Value interface {
	Type() string
	marshalMessage() (proto.Message, error)
}

type ValueOut interface {
	Value
	unmarshalAny(any *anypb.Any) error
}

func Interface(value *pb.Value) (interface{}, error) {
	return unmarshalAny(value.Value)
}

func Unmarshal(value *pb.Value, out ValueOut) error {
	return out.unmarshalAny(value.Value)
}

func Marshal(value Value) (*pb.Value, error) {
	message, err := value.marshalMessage()
	if err != nil {
		return nil, err
	}
	any, err := anypb.New(message)
	if err != nil {
		return nil, err
	}
	return &pb.Value{
		Value: any,
	}, nil
}

func anyAs(any *anypb.Any, message proto.Message) bool {
	if any.MessageIs(message) {
		if err := any.UnmarshalTo(message); err != nil {
			panic(fmt.Errorf("could not unmarshal message: %w", err))
		}
		return true
	}
	return false
}

func marshalInterface(v interface{}) (*anypb.Any, error) {
	var value Value
	switch t := v.(type) {
	case bool:
		value = Bool(t)
	case int64:
		value = Int64(t)
	case int32:
		value = Int64(t)
	case int16:
		value = Int64(t)
	case int8:
		value = Int64(t)
	case int:
		value = Int64(t)
	case float32:
		value = Float64(t)
	case float64:
		value = Float64(t)
	case string:
		value = String(t)
	case []interface{}:
		value = List(t)
	case time.Time:
		value = Time{Time: t}
	case map[string]interface{}:
		value = Map(t)
	default:
		value = JSON{
			Class: "java.lang.Object",
			Value: v,
		}
	}
	msg, err := value.marshalMessage()
	if err != nil {
		return nil, err
	}
	return anypb.New(msg)
}

func unmarshalAny(any *anypb.Any) (interface{}, error) {
	var asBool pb.BoolValue
	if anyAs(any, &asBool) {
		return asBool.Value, nil
	}
	var asString pb.StringValue
	if anyAs(any, &asString) {
		return asString.Value, nil
	}
	var asInt pb.IntValue
	if anyAs(any, &asInt) {
		return asInt.Value, nil
	}
	var asDouble pb.DoubleValue
	if anyAs(any, &asDouble) {
		return asDouble.Value, nil
	}
	var asTime pb.TimeValue
	if anyAs(any, &asTime) {
		return asTime.Value, nil
	}
	var asList pb.ListValue
	if anyAs(any, &asList) {
		var list []interface{}
		for _, v := range asList.Value {
			value, err := unmarshalAny(v.Value)
			if err != nil {
				return nil, err
			}
			list = append(list, value)
		}
		return list, nil
	}
	var asMap pb.MapValue
	if anyAs(any, &asMap) {
		toMap := make(map[string]interface{})
		for key, v := range asMap.Value {
			value, err := unmarshalAny(v.Value)
			if err != nil {
				return nil, err
			}
			toMap[key] = value
		}
		return toMap, nil
	}
	var asJSON pb.JSONValue
	if anyAs(any, &asJSON) {
		var v interface{}
		if err := json.Unmarshal(asJSON.Json, &v); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSONValue: %v", err)
		}
		return v, nil
	}
	return nil, fmt.Errorf("unsupported type: %s", any.TypeUrl)
}
