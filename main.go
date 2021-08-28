// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"go.justen.tech/goodwill/internal/command"
	"os"
)

//go:generate protoc -I ./src/main/proto --go_out=. --go_opt=module=go.justen.tech/goodwill --go-grpc_out=. --go-grpc_opt=module=go.justen.tech/goodwill ./src/main/proto/config.proto ./src/main/proto/secret.proto ./src/main/proto/docker.proto ./src/main/proto/context.proto ./src/main/proto/lock.proto

var (
	Version   string
	GitCommit string
	BuildTime string
)

func main() {
	os.Exit(command.Run(command.VersionInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
	}))
}
