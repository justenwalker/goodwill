// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package value

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"go.justen.tech/goodwill/internal/pb"
)

type Bool bool

func (v Bool) Type() string {
	return "bool"
}

func (v *Bool) unmarshalAny(any *anypb.Any) error {
	var message pb.BoolValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	*v = Bool(message.Value)
	return nil
}

func (v Bool) marshalMessage() (proto.Message, error) {
	return &pb.BoolValue{
		Value: bool(v),
	}, nil
}
