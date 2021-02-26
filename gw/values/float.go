// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Float64 float64

func (v Float64) Type() string {
	return "float64"
}

func (v *Float64) unmarshalAny(any *anypb.Any) error {
	var message pb.DoubleValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	*v = Float64(message.Value)
	return nil
}

func (v Float64) marshalMessage() (proto.Message, error) {
	return &pb.DoubleValue{
		Value: float64(v),
	}, nil
}
