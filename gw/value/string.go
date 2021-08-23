// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package value

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type String string

func (v String) Type() string {
	return "string"
}

func (v *String) unmarshalAny(any *anypb.Any) error {
	var message pb.StringValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	*v = String(message.Value)
	return nil
}

func (v String) marshalMessage() (proto.Message, error) {
	return &pb.StringValue{
		Value: string(v),
	}, nil
}
