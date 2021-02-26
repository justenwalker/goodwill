// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type List []interface{}

func (v List) Type() string {
	return "list"
}

func (v *List) unmarshalAny(any *anypb.Any) error {
	var message pb.ListValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	var result []interface{}
	for _, v := range message.Value {
		value, err := unmarshalAny(v.Value)
		if err != nil {
			return err
		}
		result = append(result, value)
	}
	*v = result
	return nil
}

func (v List) marshalMessage() (proto.Message, error) {
	var list []*pb.Value
	for _, val := range v {
		any, err := marshalInterface(val)
		if err != nil {
			return nil, err
		}
		list = append(list, &pb.Value{Value: any})
	}
	return &pb.ListValue{
		Value: list,
	}, nil
}
