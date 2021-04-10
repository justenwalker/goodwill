// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"fmt"
	"go.justen.tech/goodwill/internal/parse"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Run() int {
	var flg flags
	if !flg.Parse() {
		return 128
	}
	if flg.Debug {
		SetDebug(log.New(os.Stderr, "[DEBUG] ", 0))
	}
	data, err := parsePackage(flg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := compile(flg, *data); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

type ParsedData struct {
	Functions []parse.TaskFunction
}

func parsePackage(flg flags) (*ParsedData, error) {
	files, err := getGoFiles(flg.Dir, flg)
	if err != nil {
		return nil, fmt.Errorf("could not get go files to inspect: %w", err)
	}
	info, err := parse.Package(flg.Dir, files)
	if err != nil {
		return nil, err
	}
	if info.AstPkg.Name != "main" {
		return nil, fmt.Errorf("package must be named 'main', was %s", info.AstPkg.Name)
	}
	funcs, err := parse.TaskFunctions(flg.Dir, info)
	if err != nil {
		return nil, err
	}
	return &ParsedData{
		Functions: funcs,
	}, nil
}

func compile(flg flags, td ParsedData) error {
	file, err := mainFile(flg.Dir, td)
	if err != nil {
		return fmt.Errorf("could not generate main file: %w", err)
	}
	defer os.Remove(file)
	files, err := getGoFiles(flg.Dir, flg)
	if err != nil {
		return fmt.Errorf("could not get go files to compile: %w", err)
	}
	if err := compileGoFiles(files, flg, flg.Output); err != nil {
		return fmt.Errorf("failed to compile goodwill binary: %w", err)
	}
	return nil
}

func mainFile(dir string, data ParsedData) (string, error) {
	f, err := ioutil.TempFile(dir, "goodwill-main.*.go")
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), renderMain(f, data)
}

func compileGoFiles(files []string, flg flags, out string) error {
	args := []string{"build"}
	args = append(args, flg.GoFlags...)
	args = append(args, "-tags=goodwill", "-o", out)
	args = append(args, files...)
	cmd := exec.Command(flg.GoCmd, args...)
	debug.Println("running", cmd.Path, strings.Join(cmd.Args[1:], " "))
	debug.Printf("GOOS=%s, GOARCH=%s", flg.OS, flg.Arch)
	cmd.Env = GoEnv(flg.OS, flg.Arch)
	cmd.Dir = flg.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getGoFiles(cwd string, flg flags) ([]string, error) {
	args := []string{"list"}
	args = append(args, flg.GoFlags...)
	args = append(args, "-e", "-f", `{{ join .GoFiles "||" }}`)
	debug.Println("getting all non-goodwill files in", cwd)
	str, err := RunEnvCommand(GoEnv(flg.OS, flg.Arch), cwd, flg.GoCmd, args...)
	if err != nil {
		return nil, err
	}
	debug.Println("non-goodwill files in", str)
	exclude := make(map[string]struct{})
	for _, f := range strings.Split(str, "||") {
		exclude[f] = struct{}{}
	}
	args = append(args, "-tags=goodwill")
	str, err = RunEnvCommand(GoEnv(flg.OS, flg.Arch), cwd, flg.GoCmd, args...)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, f := range strings.Split(str, "||") {
		if _, ok := exclude[f]; !ok {
			files = append(files, f)
		}
	}
	debug.Println("goodwill files:", strings.Join(files, ""))
	return files, nil
}
