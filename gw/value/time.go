// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package value

import (
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"go.justen.tech/goodwill/internal/pb"
)

type Time struct {
	time.Time
}

func (v Time) Type() string {
	return "time"
}

func (v *Time) unmarshalAny(any *anypb.Any) error {
	var intValue pb.IntValue
	if anyAs(any, &intValue) {
		*v = Time{
			Time: time.Unix(0, intValue.Value*int64(time.Millisecond)),
		}
		return nil
	}
	var message pb.TimeValue
	if err := any.UnmarshalTo(&message); err != nil {
		return err
	}
	*v = Time{
		Time: time.Unix(0, message.Value*int64(time.Millisecond)),
	}
	return nil
}

func (v Time) marshalMessage() (proto.Message, error) {
	return &pb.TimeValue{
		Value: v.UnixNano() / int64(time.Millisecond),
	}, nil
}
