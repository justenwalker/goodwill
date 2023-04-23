// Copyright 2021, Justen Walker and the goodwill contributors
// SPDX-License-Identifier: Apache-2.0

package mage

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ZipFile struct {
	SourceReader io.Reader
	Source       string
	Dest         string
}

func WriteZip(w io.Writer, files []ZipFile) error {
	zw := zip.NewWriter(w)
	for _, file := range files {
		if file.SourceReader != nil {
			if err := addZipReader(zw, file.SourceReader, file.Dest); err != nil {
				return err
			}
			continue
		}
		stat, err := os.Stat(file.Source)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			err = addZipDir(zw, file.Source, file.Dest)
		} else {
			err = addZipFile(zw, file.Source, file.Dest)
		}
		if err != nil {
			return err
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("close zip file: %w", err)
	}
	return nil
}

func addZipDir(zw *zip.Writer, dir string, dest string) error {
	return filepath.Walk(dir, func(filename string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return addZipFile(zw, filename, path.Join(dest, strings.TrimPrefix(filename, dir)))
	})
}

func addZipReader(zw *zip.Writer, r io.Reader, dest string) error {
	//debug.Printf("> payload: %s <- %s", dest, filename)
	w, err := zw.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

func addZipFile(zw *zip.Writer, filename string, dest string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return addZipReader(zw, f, dest)
}
