# Goodwill

Goodwill is a plugin for [Walmart Labs Concord](https://concord.walmartlabs.com/). It provides a way to write tasks on
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
)

// Each task is a function with the following signature
// The function name should be capitalized
func MyTask(ts *gw.Task) error {
	_ = ts.Log(context.TODO(), "Hello, Goodwill!")
}
```

In your `concord.yml`, you can call the go code using the `goodwill` task.

```yaml
flows:
  default:
    - task: goodwill
      in:
        flow: mytask # case insensitive
```

