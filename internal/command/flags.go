// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

type flags struct {
	GoCmd   string
	OS      string
	Arch    string
	GoFlags []string
	goFlags string
	Debug   bool
	Dir     string
	Output  string
	Version bool
}

func (f *flags) Parse() bool {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.BoolVar(&f.Version, "version", false, "print version information")
	fs.StringVar(&f.Output, "out", "goodwill.tasks", "the output binary")
	fs.StringVar(&f.Dir, "dir", ".", "the directory containing the goodwill source")
	fs.StringVar(&f.OS, "os", runtime.GOOS, "set the GOOS")
	fs.StringVar(&f.Arch, "arch", runtime.GOARCH, "set the GOARCH")
	fs.BoolVar(&f.Debug, "debug", false, "enable debug logging")
	fs.StringVar(&f.GoCmd, "gobin", "go", "path to the go binary")
	fs.StringVar(&f.goFlags, "goflags", "", "additional flags to pass to go")
	if err := fs.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		return false
	}
	f.GoFlags = splitArgs(f.goFlags)
	return true
}

func splitArgs(str string) []string {
	var args []string
	r := strings.NewReader(str)
	var quote rune
	var esc bool
	var arglen int
	var lastQuote rune
	var buf strings.Builder
	for {
		ch, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}
		if esc {
			esc = false
			buf.WriteRune(ch)
			arglen++
			continue
		}
		if quote > 0 {
			if ch == '\\' {
				esc = true
				continue
			}
			if ch == quote {
				lastQuote = quote
				quote = 0
				continue
			}
			buf.WriteRune(ch)
			arglen++
			continue
		}
		switch ch {
		case '"', '\'':
			if ch == lastQuote {
				buf.WriteRune(ch)
				arglen++
			}
			quote = ch
		case '\\':
			esc = true
		case ',':
			quote = 0
			if arglen > 0 {
				args = append(args, buf.String())
			}
			arglen = 0
			buf.Reset()
		default:
			quote = 0
			arglen++
			buf.WriteRune(ch)
		}
		lastQuote = 0
	}
	if arglen > 0 {
		args = append(args, buf.String())
	}
	return args
}
