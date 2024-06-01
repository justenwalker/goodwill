// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

//go:build goodwill
// +build goodwill

package main

import (
	"context"

	// This tests that aliases work when parsing
	gw2 "go.justen.tech/goodwill/gw"
	val "go.justen.tech/goodwill/gw/value"
)

// ContextFunc takes a context and a task and returns an error
func ContextFunc(ctx context.Context, ts *gw2.Task) error {
	return nil
}

// ContextOutFunc takes a context and a task and returns output variables and an error
func ContextOutFunc(ctx context.Context, ts *gw2.Task) (map[string]val.Value, error) {
	return nil, nil
}
