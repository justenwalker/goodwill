// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Int64 int64

func (v Int64)  Type() string {
	return "int64"
}

func (v *Int64) unmarshalAny(any *anypb.Any) error {
	var message pb.IntValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	*v = Int64(message.Value)
	return nil
}

func (v Int64) marshalMessage() (proto.Message, error) {
	return &pb.IntValue{
		Value: int64(v),
	}, nil
}
