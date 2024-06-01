// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package lock

import (
	"context"

	"google.golang.org/grpc"

	"go.justen.tech/goodwill/internal/pb"
)

// Service manages project locks
type Service struct {
	client pb.LockServiceClient
}

// NewService creates a new service to manage project locks
func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		client: pb.NewLockServiceClient(conn),
	}
}

// ProjectLock requests an exclusive lock on the given name in the current project
func (c *Service) ProjectLock(ctx context.Context, name string) error {
	_, err := c.client.ProjectLock(ctx, &pb.Lock{Name: name})
	return err
}

// ProjectUnlock unlocks a named lock previously locked by ProjectLock
func (c *Service) ProjectUnlock(ctx context.Context, name string) error {
	_, err := c.client.ProjectUnlock(ctx, &pb.Lock{Name: name})
	return err
}
