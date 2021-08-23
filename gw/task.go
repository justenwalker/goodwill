// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package gw

import (
	"context"
	"fmt"
	"go.justen.tech/goodwill/gw/config"
	"go.justen.tech/goodwill/gw/docker"
	"go.justen.tech/goodwill/gw/jsonstore"
	"go.justen.tech/goodwill/gw/kv"
	"go.justen.tech/goodwill/gw/lock"
	"go.justen.tech/goodwill/gw/secret"
	"go.justen.tech/goodwill/gw/taskcontext"
	"go.justen.tech/goodwill/gw/value"
	"google.golang.org/grpc"
)

// newTask creates a new Task connection between the client and the GRPC service on the agent
func newTask(addr string, rt runtime, opts ...grpc.DialOption) (*Task, error) {
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to grpc server %q: %w", addr, err)
	}
	return &Task{
		OrgName:    rt.OrgName,
		WorkingDir: rt.WorkingDir,
		ProcessID:  rt.ProcessID,
		conn:       conn,
	}, nil
}

// runtime contains the process runtime information
type runtime struct {
	OrgName    string
	WorkingDir string
	ProcessID  string
}

type Task struct {
	OrgName    string
	WorkingDir string
	ProcessID  string
	conn       *grpc.ClientConn
}

type TaskRunner interface {
	Run(ctx context.Context, ts *Task) (map[string]value.Value, error)
}

type TaskRunnerFunc func(ctx context.Context, ts *Task) (map[string]value.Value, error)

func (f TaskRunnerFunc) Run(ctx context.Context, ts *Task) (map[string]value.Value, error) {
	return f(ctx, ts)
}

const (
	logLineVar      = "_goodwill_log_line"
	logLineCallExpr = `${log.call(` + logLineVar + `)}`
)

// Log emits a log entry on the process.
func (c *Task) Log(ctx context.Context, format string, v ...interface{}) error {
	return c.Context().EvaluateParams(ctx, logLineCallExpr, value.Discard, map[string]value.Value{
		logLineVar: value.String(fmt.Sprintf(format, v...)),
	})
}

// Docker returns a service for running Docker containers
func (c *Task) Docker() *docker.Service {
	return docker.NewService(c.conn)
}

// Context returns a service for interacting with the process context
// such as setting and getting variabes, and evaluating the results of Java Expressions in JEL 3.0
func (c *Task) Context() *taskcontext.Service {
	return taskcontext.NewService(c.conn)
}

// Config returns a service for retrieving process configuration
func (c *Task) Config() *config.Service {
	return config.NewService(c.conn)
}

// Secret returns a service for manipulating secrets
// within the current Organization
func (c *Task) Secret() *secret.Service {
	return secret.NewService(c.OrgName, c.conn)
}

// JSONStore returns a jsonstore.Store for the given store name
// The store should already exist, otherwise operations will fail.
// This store should be in the current Org of the process.
func (c *Task) JSONStore(name string) *jsonstore.Store {
	return jsonstore.NewService(c.OrgName, name, c.conn)
}

// Lock returns a service for setting project-level locks
func (c *Task) Lock() *lock.Service {
	return lock.NewService(c.conn)
}

// KV returns a service for interacting with the project's key-value storage system
func (c *Task) KV() *kv.Service {
	return kv.NewService(c.conn)
}

func (c *Task) Close() error {
	return c.conn.Close()
}
