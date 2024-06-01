// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"go.justen.tech/goodwill/internal/command"
)

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
