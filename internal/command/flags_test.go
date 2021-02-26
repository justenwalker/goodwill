// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"reflect"
	"testing"
)

func TestSplitFlags(t *testing.T) {
	tests := []struct {
		name     string
		parse    string
		expected []string
	}{
		{"one", "foo", []string{"foo"}},
		{"two", "foo,bar", []string{"foo", "bar"}},
		{"three", "foo,bar,baz", []string{"foo", "bar", "baz"}},
		{"comma-1", "foo,bar,", []string{"foo", "bar"}},
		{"comma-2", ",foo,bar,", []string{"foo", "bar"}},
		{"double-quote", `"foo"`, []string{"foo"}},
		{"double-quote-comma", `"foo,bar"`, []string{"foo,bar"}},
		{"single-quote", `'foo'`, []string{"foo"}},
		{"single-quote-comma", `'foo,bar'`, []string{"foo,bar"}},
		{"embed-double-quote", `"foo""bar",baz`, []string{`foo"bar`, "baz"}},
		{"embed-single-quote", `'foo''bar',baz`, []string{`foo'bar`, "baz"}},
		{"no-embed-double-quote", `'foo""bar',baz`, []string{`foo""bar`, "baz"}},
		{"no-embed-single-quote", `"foo''bar",baz`, []string{`foo''bar`, "baz"}},
		{"esc-single-quote", `\'foo\'`, []string{"'foo'"}},
		{"esc-comma", `foo\,bar`, []string{"foo,bar"}},
		{"esc-double-quote", `"foo\"bar",baz`, []string{`foo"bar`, "baz"}},
		{"esc-single-quote", `'foo\'bar',baz`, []string{`foo'bar`, "baz"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := splitArgs(test.parse)
			if !reflect.DeepEqual(args, test.expected) {
				t.Fatalf("%s != %s", args, test.expected)
			}
		})
	}
}
