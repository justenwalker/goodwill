// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package docker

import (
	"context"
	"fmt"
	"go.justen.tech/goodwill/internal/pb"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
)

type Service struct {
	client pb.DockerServiceClient
}

func NewService(conn grpc.ClientConnInterface) *Service {
	return &Service{
		client: pb.NewDockerServiceClient(conn),
	}
}

// RunContainer runs a docker container. The container being run should exit.
// If the container exits with a non-zero exit code, an error implementing `interface { ExitCode() int }` will be returned
// from which the exit code may be extracted.
func (c *Service) RunContainer(ctx context.Context, image string, opts ...Option) error {
	options := Options{
		Image:  image,
		Stdout: ioutil.Discard,
		Stderr: ioutil.Discard,
	}
	for _, opt := range opts {
		opt.Apply(&options)
	}
	stream, err := c.client.RunContainer(ctx, &pb.DockerContainerSpec{
		Image:            options.Image,
		Name:             options.Name,
		User:             options.User,
		WorkDir:          options.WorkDir,
		EntryPoint:       options.EntryPoint,
		Command:          options.Command,
		Env:              options.Env,
		EnvFile:          options.EnvFile,
		Labels:           options.Labels,
		ForcePull:        options.ForcePull,
		Hosts:            options.Hosts,
		StdoutFilePath:   options.StdoutFilePath,
		RedirectStdError: options.RedirectStdError,
	})
	stderr := options.Stderr
	if stderr == nil {
		stderr = ioutil.Discard
	}
	stdout := options.Stdout
	if stdout == nil {
		stdout = ioutil.Discard
	}
	if err != nil {
		return err
	}
	var status int
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("docker: run container error %w", err)
		}
		result := m.GetResult()
		switch t := result.(type) {
		case *pb.DockerContainerResult_Status:
			status = int(t.Status)
		case *pb.DockerContainerResult_Stderr:
			stderr.Write([]byte(t.Stderr))
			stderr.Write([]byte{'\n'})
			if options.StderrCallback != nil {
				options.StderrCallback(t.Stderr)
			}
		case *pb.DockerContainerResult_Stdout:
			stdout.Write([]byte(t.Stdout))
			stdout.Write([]byte{'\n'})
			if options.StdoutCallback != nil {
				options.StdoutCallback(t.Stdout)
			}
		}
	}
	if status != 0 {
		return exitCodeErr(status)
	}
	return nil
}

type exitCodeErr int

func (e exitCodeErr) Error() string {
	return fmt.Sprintf("container exited with status: %d", int(e))
}

func (e exitCodeErr) ExitCode() int {
	return int(e)
}
