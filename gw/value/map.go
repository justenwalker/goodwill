// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package value

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Map map[string]interface{}

func (v Map) Type() string {
	return "map"
}

func (v *Map) unmarshalAny(any *anypb.Any) error {
	var message pb.MapValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	result := make(map[string]interface{})
	for key, any := range message.Value {
		value, err := unmarshalAny(any.Value)
		if err != nil {
			return err
		}
		result[key] = value
	}
	*v = result
	return nil
}

func (v Map) marshalMessage() (proto.Message, error) {
	mapAny := make(map[string]*pb.Value)
	for key, val := range v {
		any, err := marshalInterface(val)
		if err != nil {
			return nil, err
		}
		mapAny[key] = &pb.Value{
			Value: any,
		}
	}
	return &pb.MapValue{
		Value: mapAny,
	}, nil
}
