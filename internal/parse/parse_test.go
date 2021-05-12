// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package parse

import (
	"sort"
	"testing"
)

func TestParse(t *testing.T) {
	info, err := Package("./testdata", []string{"goodwill.go", "goodwill2.go"})
	if err != nil {
		t.Fatal(err)
	}
	expected := []TaskFunction{
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
	}
	sort.Sort(ByName(expected))
	ffs, err := TaskFunctions("./testdata", info)
	if err != nil {
		t.Fatal(err)
	}
	sort.Sort(ByName(ffs))
	for i, ff := range ffs {
		if expected[i] != ff {
			t.Errorf("Expected Func %d to be %+v, got %+v", i, expected[i], ff)
		}
	}
}
