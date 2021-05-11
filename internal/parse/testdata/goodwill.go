// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// +build goodwill

package main

import (
	"context"
	"fmt"
	"go.justen.tech/goodwill/gw"
	"go.justen.tech/goodwill/gw/value"
)

// PrintWorkDir prints the working directory
func PrintWorkDir(ts *gw.Task) error {
	tc := ts.Context()
	var value value.String
	if err := tc.Evaluate(context.TODO(), "${workDir}", &value); err != nil {
		return err
	}
	fmt.Println("workDir:", value)
	return nil
}
