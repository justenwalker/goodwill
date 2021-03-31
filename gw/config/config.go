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

// NewService creates a new config.Service for getting task configuration
func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		client: pb.NewConfigServiceClient(conn),
	}
}

// ProcessConfig is the configuration of the running process
type ProcessConfig struct {
	// ID is the UUID of the process
	ID string
	// WorkingDirectory is the current working directory of the process
	WorkingDirectory string
	// ProjectInfo has the information about the Project in which the process is running (optional)
	ProjectInfo ProjectInfo
	// RepositoryInfo has information about to the Repository inside the Project that spawned the process (optional)
	RepositoryInfo RepositoryInfo
	// APIConfiguration contains connection details to the Concord API
	APIConfiguration APIConfiguration
}

// ProjectInfo contains information about the project in which the process is running
type ProjectInfo struct {
	// OrgID is the UUID of the Organization
	OrgID string
	// OrgName is the name of the Organization
	OrgName string
	// ProjectID is the UUID of the Project
	ProjectID string
	// ProjectName is the name of the Project
	ProjectName string
}

// RepositoryInfo contains details about the repository that is configured for the process
type RepositoryInfo struct {
	// ID is the UUID of the Repository
	ID string
	// Name is the name of the Repository
	Name string
	// URL is the URL of the Repository
	URL string
}

// APIConfiguration has the api details and credentials for communicating directly with the Concord API
// as part of the process
type APIConfiguration struct {
	BaseURL      string
	SessionToken string
}

// Configuration returns the task configuration
func (c *Service) Configuration(ctx context.Context) (*ProcessConfig, error) {
	resp, err := c.client.GetConfiguration(ctx, &pb.ConfigurationRequest{})
	if err != nil {
		return nil, err
	}
	return &ProcessConfig{
		ID:               resp.ProcessID,
		WorkingDirectory: resp.WorkingDirectory,
		ProjectInfo: ProjectInfo{
			OrgID:       resp.ProjectInfo.OrgID,
			OrgName:     resp.ProjectInfo.OrgName,
			ProjectID:   resp.ProjectInfo.ProjectID,
			ProjectName: resp.ProjectInfo.ProjectName,
		},
		RepositoryInfo: RepositoryInfo{
			ID:   resp.RepositoryInfo.RepoID,
			Name: resp.RepositoryInfo.RepoName,
			URL:  resp.RepositoryInfo.RepoURL,
		},
		APIConfiguration: APIConfiguration{
			BaseURL:      resp.ApiConfiguration.BaseURL,
			SessionToken: resp.ApiConfiguration.SessionToken,
		},
	}, nil
}
