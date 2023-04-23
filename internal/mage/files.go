// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package mage

import (
	"embed"
	"io"
	"os"
)

//go:embed files/*
var Files embed.FS

func WriteFile(filename string, wfn func(w io.Writer) error) error {
	out, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(0o644))
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	return wfn(out)
}
