// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package mage

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerateConcordYaml(t *testing.T) {
	tests := []struct {
		name   string
		params ConcordParams
	}{
		{"default", ConcordParams{}},
		{"dependencies", ConcordParams{Dependencies: true}},
		{"dependencies-version", ConcordParams{Dependencies: true, Version: "0.4.1"}},
		{"use-docker", ConcordParams{UseDocker: true}},
		{"go-version", ConcordParams{GoVersion: "1.16.4"}},

		{"default-v2", ConcordParams{Runtime: ConcordRuntimeV2}},
		{"dependencies-v2", ConcordParams{Runtime: ConcordRuntimeV2, Dependencies: true}},
		{"dependencies-version-v2", ConcordParams{Runtime: ConcordRuntimeV2, Dependencies: true, Version: "0.4.1"}},
		{"use-docker-v2", ConcordParams{Runtime: ConcordRuntimeV2, UseDocker: true}},
		{"go-version-v2", ConcordParams{Runtime: ConcordRuntimeV2, GoVersion: "1.16.4"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := goldie.New(t)
			var buf bytes.Buffer
			require.NoError(t, GenerateConcordYaml(&buf, test.params))
			g.Assert(t, test.name, buf.Bytes())
		})
	}
}
