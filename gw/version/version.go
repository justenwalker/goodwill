// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"

	"github.com/masterminds/semver"
)

func CheckCompatibility(serverVersion string, minVersion string) error {
	sv, err := semver.NewVersion(serverVersion)
	if err != nil {
		return fmt.Errorf("could not parse server version %q: %w", serverVersion, err)
	}
	mv, err := semver.NewVersion(minVersion)
	if err != nil {
		panic(fmt.Errorf("could not parse minimum version %q: %w", minVersion, err))
	}
	if sv.Compare(mv) < 0 {
		return fmt.Errorf("server version %q is older than minimum version required %q", serverVersion, minVersion)
	}
	return nil
}
