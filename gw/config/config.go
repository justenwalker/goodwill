// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"context"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/grpc"
)

// Service gets task configuration
type Service struct {
	client pb.ConfigServiceClient
}

func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		client: pb.NewConfigServiceClient(conn),
	}
}

type TaskConfig struct {
	ProcessID        string
	WorkingDirectory string
	ProjectInfo      ProjectInfo
	RepositoryInfo   RepositoryInfo
	APIConfiguration APIConfiguration
}

type ProjectInfo struct {
	OrgID       string
	OrgName     string
	ProjectID   string
	ProjectName string
}

type RepositoryInfo struct {
	RepoID   string
	RepoName string
	RepoURL  string
}

type APIConfiguration struct {
	BaseURL      string
	SessionToken string
}

// Configuration returns the task configuration
func (c *Service) Configuration(ctx context.Context) (*TaskConfig, error) {
	resp, err := c.client.GetConfiguration(ctx, &pb.ConfigurationRequest{})
	if err != nil {
		return nil, err
	}
	return &TaskConfig{
		ProcessID:        resp.ProcessID,
		WorkingDirectory: resp.WorkingDirectory,
		ProjectInfo: ProjectInfo{
			OrgID:       resp.ProjectInfo.OrgID,
			OrgName:     resp.ProjectInfo.OrgName,
			ProjectID:   resp.ProjectInfo.ProjectID,
			ProjectName: resp.ProjectInfo.ProjectName,
		},
		RepositoryInfo: RepositoryInfo{
			RepoID:   resp.RepositoryInfo.RepoID,
			RepoName: resp.RepositoryInfo.RepoName,
			RepoURL:  resp.RepositoryInfo.RepoURL,
		},
		APIConfiguration: APIConfiguration{
			BaseURL:      resp.ApiConfiguration.BaseURL,
			SessionToken: resp.ApiConfiguration.SessionToken,
		},
	}, nil
}
