# Goodwill

Goodwill is a plugin for [Walmart Labs Concord](https://concord.walmartlabs.com/). It provides a way to write tasks in
Go by providing a bridge to the Java runtime over GRPC.

## Quick Start

Goodwill tasks are written in go files which have a build tag `goodwill`

In your project root, create a file named something like `goodwill.go`.


```go
// +build goodwill

package main

import (
	"context"

	"go.justen.tech/goodwill/gw"
	"go.justen.tech/goodwill/gw/value"
)

// Each task is a function with the following signature
// The function name should be capitalized

// Both the ctx context.Context arg and the map[string]value.Value return are optional
// and may be omitted from the method signature

// Default Task, Says Hello
func Default(ctx context.Context, ts *gw.Task) (map[string]value.Value, error) {
	_ = ts.Log(ctx, "Hello, Goodwill!")
	return nil, nil
}
```

In your `concord.yml`, you can call the go code using the `goodwill` task.

```yaml
configuration:
  dependencies:
    - mvn://tech.justen.concord:goodwill:0.7.1

flows:
  default:
    - task: goodwill
```


## Precompile Task Binary

By default, goodwill installs Go and compiles the tasks on the Concord server.
While the cost of this is amortized over the number of tasks you execute within a single flow,
this added time and compilation on the Concord agent may be undesirable.
In these cases, you can run `goodwill` yourself to pre-compile the tasks.

```shell
$ goodwill -os linux -arch amd64
```

If a `goodwill.tasks` binary exists in your concord payload, this will be used instead of compiling from your tasks source.

## Task Parameters

The `goodwill` task can take several optional `in` parameters.

### Common
- `task`: The name of the task to run. (default: `Default`)
- `binary`: The path to the pre-compiled goodwill tasks. (default: `goodwill.tasks`)
- `debug`: Enable debug logging for compilation. (default: `false`)

### Build
- `goVersion`: Set the version of Go to install. (Default: `1.22.3`)
- `goDockerImage`: Override the image to use when building a goodwill flow in Docker (Default: `golang:${goVersion}`)
- `useDocker`: Use a docker image to compile the goodwill binary.
- `installGo`: If compilation is required and go is not found, install it in the task workspace (Default: `true`)
- `buildDir`: Override the output directory for goodwill generated files. (default: `.goodwill`)'

### Go Environment Settings
- `GOOS`: Override the OS target of the concord agent instead of auto-detecting it.
- `GOARCH`: Override the OS Architecture of the target instead of auto-detecting it.
- `GOPROXY`: Set the `GOPROXY` environment variable
- `GONOPROXY`: Set the `GONOPROXY` environment variable
- `GOPRIVATE`: Set the `GOPRIVATE` environment variable
- `GOSUMDB`: Set the `GOSUMDB` environment variable
- `GONOSUMDB`: Set the `GONOSUMDB` environment variable

## Building and Testing

This project uses [Mage](https://magefile.org/) to produce build artifacts and to run tests.
Before building, you'll need to install Mage yourself.

```shell
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

### Building Artifacts

To build just the binary for your platform:

```shell
mage build
```

The binary will be in `./bin`.

To build and package release files, including the Concord plugin, run:

```shell
mage package
```

The files will be placed in `./dist`

### Running Tests

To run End-to-End tests, you need to install [OpenTofu](https://opentofu.org/docs/intro/install/) and [Docker](https://www.docker.com/).

```shell
mage e2e:test
```

This will bring a local Concord instance up in Docker, and make it available on `http://localhost:8001`
and run a few flows against the instance using Goodwill.
