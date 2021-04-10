// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"github.com/sebdah/goldie/v2"
	"go.justen.tech/goodwill/internal/parse"
	"testing"
)

func TestRenderMain(t *testing.T) {
	var buf bytes.Buffer
	g := goldie.New(t)
	err := renderMain(&buf, ParsedData{
		Functions: []parse.TaskFunction{
			{Name: "Default", Doc: "Default Doc"},
			{Name: "Task1", Doc: "Task 1 Doc"},
			{Name: "Task2", Doc: "Task 2 Doc"},
		},
	})
	if err != nil {
		t.Fatal("unexpected error in renderMain:", err)
	}
	g.Assert(t, "main", buf.Bytes())
}
