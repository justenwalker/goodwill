// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var debug = log.New(ioutil.Discard, "", 0)

// SetDebug sets the debug logger
func SetDebug(l *log.Logger) {
	debug = l
}

// GetEnv gets a copy of the process environment variables
// and overrides them with the environment variable map
func GetEnv(envs map[string]string) []string {
	m := environToMap(os.Environ())
	for k, v := range envs {
		m[k] = v
	}
	return mapToEnviron(m)
}

// RunCommand runs the given command
func RunCommand(name string, args ...string) error {
	out, err := RunCommandOut(name, args...)
	if err != nil {
		return err
	}
	debug.Println("stdout:\n", out)
	return nil
}

// RunCommand runs the given command and returns stdout
func RunCommandOut(name string, args ...string) (string, error) {
	return RunEnvCommand(nil, "", name, args...)
}

// RunEnvCommand runs the given command and returns stdout
func RunEnvCommand(env []string, dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stderr, stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if dir != "" {
		cmd.Dir = dir
	}
	if env != nil {
		cmd.Env = env
	} else {
		cmd.Env = CurrentGoEnv()
	}
	debug.Println("running", cmd, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		debug.Print("error running '", cmd, strings.Join(args, " "), "'",
			"\nerror:", err,
			"\nstderr:\n", stderr.String())
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

// CurrentGoEnv returns the environment with the given GOOS/GOARCH set
// to the current runtime.GOOS/GOARCH
func CurrentGoEnv() []string {
	return GoEnv(runtime.GOOS, runtime.GOARCH)
}

// GoEnv returns the environment with the given GOOS/GOARCH set
func GoEnv(os, arch string) []string {
	return GetEnv(map[string]string{
		"GOOS":   os,
		"GOARCH": arch,
	})
}

func mapToEnviron(m map[string]string) []string {
	var envs []string
	for k, v := range m {
		envs = append(envs, k+"="+v)
	}
	return envs
}

func environToMap(envs []string) map[string]string {
	out := map[string]string{}
	for _, env := range envs {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			debug.Printf("badly formed environment variable: %v", env)
			continue
		}
		key, value := kv[0], kv[1]
		out[key] = value
	}
	return out
}
