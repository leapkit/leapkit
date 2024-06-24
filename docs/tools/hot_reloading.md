---
title: Hot reloading
index: 1
---

You can include the hot reloading feature in your app by using the `rebuilder` package, which enables us to reload our application in a development environment.

## How to Use

Firstly, you need to create another entry point in your app in which this package will live, for instance, `./cmd/dev/main.go`. Then, in the `main` function, let's make use of the `rebuilder.Start()` function, setting the main entry path as the first parameter:

```go
// ./cmd/dev/main.go

package main

import "github.com/leapkit/core/rebuilder"

func main() {
	rebuilder.Start("./cmd/app/main.go")
}
```

## Options

The `Start` function allows setting options to react when files are changed in your project.

### `rebuilder.WatchExtension`

This option allows us to configure the package to listen for changes in the files with the extension you specify and reload the application accordingly. For instance:

```go
rebuilder.Start("./cmd/app/main.go",
    rebuilder.WatchExtension(".go", ".html"),
)
```

### `rebuilder.ExcludePaths`

This option configures the package to prevent reloading the app when any file within any of the paths specified changes. For instance:

```go
rebuilder.Start("./cmd/app/main.go",
    rebuilder.ExcludePaths("./internal/test/")
)
```

### `rebuilder.WithRunner`

This option enables you to add runners, which are functions with the signature func(), and are executed each time the app refreshes.


```go
func myAwesomeRun() {
    fmt.Println("I'll be executed, Wooooh!")
}

// ...
rebuilder.Start("./cmd/app/main.go",
    rebuilder.WithRunner(myAwesomeRun)
)
```

## Running your app

Once you have configured the Start function, you have to run your app from the new entry:

```bash
go run ./cmd/dev
```
