// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package gw

import (
	"context"
	"fmt"

	"go.justen.tech/goodwill/gw/config"
	"go.justen.tech/goodwill/gw/docker"
	"go.justen.tech/goodwill/gw/kv"
	"go.justen.tech/goodwill/gw/lock"
	"go.justen.tech/goodwill/gw/secret"
	"go.justen.tech/goodwill/gw/taskcontext"
	"go.justen.tech/goodwill/gw/values"
	"google.golang.org/grpc"
)

func NewTask(addr string, rt Runtime, opts ...grpc.DialOption) (*Task, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to grpc server %q: %w", addr, err)
	}
	return &Task{
		Runtime: rt,
		conn:    conn,
	}, nil
}

type Runtime struct {
	OrgName    string
	WorkingDir string
	ProcessID  string
}

type Task struct {
	Runtime
	conn *grpc.ClientConn
}

const (
	logLineVar      = "_goodwill_log_line"
	logLineCallExpr = `${log.call(` + logLineVar + `)}`
)

func (c *Task) Log(ctx context.Context, format string, v ...interface{}) error {
	return c.Context().EvaluateParams(ctx, logLineCallExpr, values.Discard, map[string]values.Value{
		logLineVar: values.String(fmt.Sprintf(format, v...)),
	})
}

func (c *Task) Docker() *docker.Service {
	return docker.NewService(c.conn)
}

func (c *Task) Context() *taskcontext.Service {
	return taskcontext.NewService(c.conn)
}

func (c *Task) Config() *config.Service {
	return config.NewService(c.conn)
}

func (c *Task) Secret() *secret.Service {
	return secret.NewService(c.OrgName, c.conn)
}

func (c *Task) Lock() *lock.Service {
	return lock.NewService(c.conn)
}

func (c *Task) KV() *kv.Service {
	return kv.NewService(c.conn)
}

func (c *Task) Close() error {
	return c.conn.Close()
}
