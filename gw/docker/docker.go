// Copyright 2021, Justen Walker
// SPDX-License-Identifier: Apache-2.0

package docker

import "io"

// WorkspaceDir is the path of the task's mounted workspace inside a container
const WorkspaceDir = "/workspace"

// Options are the container run options
type Options struct {
	// Image is the docker image reference for the container. (Required)
	Image string
	// Name sets the name of the container
	Name string
	// User sets a different user to run unside the container
	User string
	// WorkDir specifies an alternate starting directory
	WorkDir string
	// EntryPoint changes the command that is the entrypoint of the docker image
	EntryPoint string
	// Command is the list of command arguments to pass to the container
	Command []string
	// Cpu is the amount of cpu to limit the container
	Cpu string
	// Memory is the amount of memory to limit the container
	Memory string
	// Env contains additional environment variables to set
	Env map[string]string
	// EnvFile loads the given file with the list of environment variables and values
	EnvFile string
	// Labels contains a list of labels to apply to the container
	Labels map[string]string
	// ForcePull forces a docker pull before running the container, useful if the tag was updated
	ForcePull bool
	// Hosts is a list of additional hosts that the container can resolve, like from /etc/hosts
	// Hosts are given `name:ip` strings like 'foo:10.0.0.3'
	Hosts []string
	// StdoutFilePath directs all standard error to a file
	StdoutFilePath string
	// RedirectStdError redirects all standard err output into standard out
	RedirectStdError bool
	// Stdout directs all standard out to the given writer
	Stdout io.Writer
	// Stderr directs all standard err to the given writer
	Stderr io.Writer
	// StdoutCallback is an optional function that is called with every output line of the container's standard out
	StdoutCallback func(string)
	// StderrCallback is an optional function that is called with every output line of the container's standard err
	StderrCallback func(string)
}

func (o Options) Apply(t *Options) {
	*t = o
}

type Option interface {
	Apply(opts *Options)
}

// WorkDir sets the working directory of the container
func WorkDir(workDir string) Option {
	return workDirOption(workDir)
}

type workDirOption string

func (o workDirOption) Apply(spec *Options) {
	spec.WorkDir = string(o)
}

// User sets the user/uid of the container process
func User(user string) Option {
	return userOption(user)
}

type userOption string

func (o userOption) Apply(spec *Options) {
	spec.User = string(o)
}

// Entrypoint overrides the container image's entrypoint
func Entrypoint(entryPoint string) Option {
	return entryPointOption(entryPoint)
}

type entryPointOption string

func (o entryPointOption) Apply(spec *Options) {
	spec.EntryPoint = string(o)
}

// Command sets the command arguments to the container image
func Command(cmd ...string) Option {
	return commandOption(cmd)
}

type commandOption []string

func (o commandOption) Apply(spec *Options) {
	spec.Command = o
}

// CPU sets the cpu shares
func CPU(cpu string) Option {
	return cpuOption(cpu)
}

type cpuOption string

func (o cpuOption) Apply(spec *Options) {
	spec.Cpu = string(o)
}

// Memory sets the maximum memory the container may use
func Memory(mem string) Option {
	return memoryOption(mem)
}

type memoryOption string

func (o memoryOption) Apply(spec *Options) {
	spec.Memory = string(o)
}

// Env sets an environment variable in the container
func Env(key string, value string) Option {
	return envOption{key: key, value: value}
}

type envOption struct {
	key   string
	value string
}

// EnvFile read environment variables from a file
func EnvFile(envFile string) Option {
	return envFileOption(envFile)
}

type envFileOption string

func (o envFileOption) Apply(spec *Options) {
	spec.EnvFile = string(o)
}

func (o envOption) Apply(spec *Options) {
	if spec.Env == nil {
		spec.Env = make(map[string]string)
	}
	spec.Env[o.key] = o.value
}

// Label sets a container label
func Label(key string, value string) Option {
	return labelOption{key: key, value: value}
}

type labelOption struct {
	key   string
	value string
}

func (o labelOption) Apply(spec *Options) {
	if spec.Labels == nil {
		spec.Labels = make(map[string]string)
	}
	spec.Labels[o.key] = o.value
}

// ForcePull sets the containers force-pull option
func ForcePull(forcePull bool) Option {
	return forcePullOption(forcePull)
}

type forcePullOption bool

func (o forcePullOption) Apply(spec *Options) {
	spec.ForcePull = bool(o)
}

// RedirectStderr redirects stderr to stdout, combining both streams into one.
func RedirectStderr(redirect bool) Option {
	return redirectStderrOption(redirect)
}

type redirectStderrOption bool

func (o redirectStderrOption) Apply(spec *Options) {
	spec.RedirectStdError = bool(o)
}

// Stdout sets the stdout stream
func Stdout(w io.Writer) Option {
	return stdoutOption{writer: w}
}

type stdoutOption struct {
	writer io.Writer
}

func (o stdoutOption) Apply(spec *Options) {
	spec.Stdout = o.writer
}

// Stderr sets the stderr stream
func Stderr(w io.Writer) Option {
	return stderrOption{writer: w}
}

type stderrOption struct {
	writer io.Writer
}

func (o stderrOption) Apply(spec *Options) {
	spec.Stderr = o.writer
}

// StdoutCallback sets the function that gets called for every stdout line
func StdoutCallback(fn func(line string)) Option {
	return stdoutCallbackOption{fn: fn}
}

type stdoutCallbackOption struct {
	fn func(string)
}

func (o stdoutCallbackOption) Apply(spec *Options) {
	spec.StdoutCallback = o.fn
}

// StderrCallback sets the function that gets called for every stderr line
func StderrCallback(fn func(line string)) Option {
	return stderrCallbackOption{fn: fn}
}

type stderrCallbackOption struct {
	fn func(string)
}

func (o stderrCallbackOption) Apply(spec *Options) {
	spec.StderrCallback = o.fn
}
