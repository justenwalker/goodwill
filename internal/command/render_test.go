// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"

	"go.justen.tech/goodwill/internal/parse"
)

func TestRenderMain(t *testing.T) {
	var buf bytes.Buffer
	g := goldie.New(t)
	err := renderMain(&buf, ParsedData{
		Functions: []parse.TaskFunction{
			{
				Name:    "ContextFunc",
				Doc:     "ContextFunc takes a context and a task and returns an error",
				Context: true,
			},
			{
				Name:    "ContextOutFunc",
				Doc:     "ContextOutFunc takes a context and a task and returns output variables and an error",
				Context: true,
				OutVars: true,
			},
			{
				Name: "Func",
				Doc:  "Func only takes a task and returns an error",
			},
			{
				Name:    "OutFunc",
				Doc:     "OutFunc takes a task and returns output variables and an error",
				OutVars: true,
			},
		},
	})
	if err != nil {
		t.Fatal("unexpected error in renderMain:", err)
	}
	g.Assert(t, "main", buf.Bytes())
}
