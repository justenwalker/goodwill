// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package value

import (
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const Discard = nullOut(0)

type nullOut int

func (n nullOut)  Type() string {
	return "null"
}

func (n nullOut) marshalMessage() (proto.Message, error) {
	return &pb.NullValue{}, nil
}

func (n nullOut) unmarshalAny(any *anypb.Any) error {
	return nil
}
