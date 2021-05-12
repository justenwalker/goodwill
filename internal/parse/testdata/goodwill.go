// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

// +build goodwill

package main

import (
	"go.justen.tech/goodwill/gw"
	"go.justen.tech/goodwill/gw/value"
)

// Func only takes a task and returns an error
func Func(ts *gw.Task) error {
	return nil
}

// OutFunc takes a task and returns output variables and an error
func OutFunc(ts *gw.Task) (map[string]value.Value, error) {
	return nil, nil
}
