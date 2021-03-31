// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"encoding/json"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
)

type JSON struct {
	Class string
	Value interface{}
}

func (v JSON) Type() string {
	return "json"
}

func (v *JSON) unmarshalAny(any *anypb.Any) error {
	val, err := unmarshalAny(any)
	if err != nil {
		return err
	}
	v.Value = val
	v.Class = v.javaClass()
	return nil
}

func (v JSON) marshalMessage() (proto.Message, error) {
	data, err := json.Marshal(v.Value)
	if err != nil {
		return nil, err
	}
	return &pb.JSONValue{
		Class: v.javaClass(),
		Json:  data,
	}, nil
}

func (v JSON) javaClass() string {
	switch v.Value.(type) {
	case bool:
		return "java.lang.Boolean"
	case int64:
		return "java.lang.Long"
	case int32:
		return "java.lang.Integer"
	case int16:
		return "java.lang.Short"
	case int8:
		return "java.lang.Byte"
	case int:
		return "java.lang.Long"
	case float32:
		return "java.lang.Float"
	case float64:
		return "java.lang.Double"
	case string:
		return "java.lang.String"
	case []interface{}:
		return "java.util.List"
	case time.Time:
		return "java.util.Date"
	case map[string]interface{}:
		return "java.util.Map"
	default:
		return "java.lang.Object"
	}
}
