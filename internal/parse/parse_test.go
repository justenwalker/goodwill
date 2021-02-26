// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package parse

import (
	"testing"
)

func TestParse(t *testing.T) {
	info, err := Package("./testdata", []string{"goodwill.go", "goodwill2.go"})
	if err != nil {
		t.Fatal(err)
	}
	ffs, err := TaskFunctions("./testdata", info)
	if err != nil {
		t.Fatal(err)
	}
	for _, ff := range ffs {
		t.Logf("> %s\n%s", ff.Doc, ff.Name)
	}
}
