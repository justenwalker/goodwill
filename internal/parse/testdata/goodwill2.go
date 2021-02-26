// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// +build goodwill

package main

import (
	c2 "context"
	"fmt"
	gw2 "go.justen.tech/goodwill/gw"
	"go.justen.tech/goodwill/gw/values"
)

// PrintProcessID prints the concord process id
func PrintProcessID(ts *gw2.Task) error {
	tc := ts.Context()
	var value values.String
	if err := tc.Evaluate(c2.TODO(), "${txId}", &value); err != nil {
		return err
	}
	fmt.Println("Process ID:", value)
	return nil
}