// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package values

import (
	"encoding/json"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type JSON struct {
	Class string
	Value interface{}
}

func (v JSON) Type() string {
	return "json"
}

func (v *JSON) unmarshalAny(any *anypb.Any) error {
	var message pb.JSONValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	if err := json.Unmarshal(message.Json, &v.Value); err != nil {
		return err
	}
	v.Class = message.Class
	return nil
}

func (v JSON) marshalMessage() (proto.Message, error) {
	data, err := json.Marshal(v.Value)
	if err != nil {
		return nil, err
	}
	return &pb.JSONValue{
		Class: v.Class,
		Json:  data,
	}, nil
}
