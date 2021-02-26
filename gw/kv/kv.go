// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package kv

import (
	"context"
	"fmt"
	"go.justen.tech/goodwill/gw/taskcontext"
	"go.justen.tech/goodwill/gw/values"
	"google.golang.org/grpc"
)

// Service to access the project key value store
type Service struct {
	ex *taskcontext.Service
}

func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		ex: taskcontext.NewService(conn),
	}
}

const (
	kvKeyVar        = "_goodwill_kv_key"
	kvValueVar      = "_goodwill_kv_val"
	kvPutStringExpr = `${kv.putString(` + kvKeyVar + `, ` + kvValueVar + `)}`
	kvPutLongExpr   = `${kv.putLong(` + kvKeyVar + `, ` + kvValueVar + `)}`
	kvGetStringExpr = `${kv.getString(` + kvKeyVar + `)}`
	kvGetLongExpr   = `${kv.getLong(` + kvKeyVar + `)}`
	kvRemoveExpr    = `${kv.remove(` + kvKeyVar + `)}`
	kvIncExpr       = `${kv.inc(` + kvKeyVar + `)}`
)

// PutString sets a string value at the given key
func (c *Service) PutString(ctx context.Context, key string, value string) error {
	return c.ex.EvaluateParams(ctx, kvPutStringExpr, values.Discard, map[string]values.Value{
		kvKeyVar:   values.String(key),
		kvValueVar: values.String(value),
	})
}

// PutLong sets an integer value at the given key
func (c *Service) PutLong(ctx context.Context, key string, value int64) error {
	return c.ex.EvaluateParams(ctx, kvPutLongExpr, values.Discard, map[string]values.Value{
		kvKeyVar:   values.String(key),
		kvValueVar: values.Int64(value),
	})
}

// GetString gets a string value at the given key
func (c *Service) GetString(ctx context.Context, key string) (string, error) {
	var o values.String
	if err := c.ex.EvaluateParams(ctx, kvGetStringExpr, &o, map[string]values.Value{
		kvKeyVar: values.String(key),
	}); err != nil {
		return "", err
	}
	return string(o), nil
}

// GetLong gets an integer value at the given key
func (c *Service) GetLong(ctx context.Context, key string) (int64, error) {
	if err := c.ex.SetVariable(ctx, kvKeyVar, values.String(key)); err != nil {
		return 0, fmt.Errorf("failed to set %s = %q", kvKeyVar, key)
	}
	var o values.Int64
	if err := c.ex.EvaluateParams(ctx, kvGetLongExpr, &o, map[string]values.Value{
		kvKeyVar: values.String(key),
	}); err != nil {
		return 0, err
	}
	return int64(o), nil
}

// Remove unsets the value at the given key
func (c *Service) Remove(ctx context.Context, key string) error {
	return c.ex.EvaluateParams(ctx, kvRemoveExpr, values.Discard, map[string]values.Value{
		kvKeyVar: values.String(key),
	})
}

// Inc increments the given key's value by 1, returning the incremented value
func (c *Service) Inc(ctx context.Context, key string) (int64, error) {
	var o values.Int64
	if err := c.ex.EvaluateParams(ctx, kvIncExpr, &o, map[string]values.Value{
		kvKeyVar: values.String(key),
	}); err != nil {
		return 0, err
	}
	return int64(o), nil
}
