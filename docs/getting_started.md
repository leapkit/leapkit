---
title: Getting Started
index: 2
---

To get started with LeapKit you should installt he `kit` CLI tool. The `kit` CLI tool is a command line interface that contains tools to accelerate your development workflow without replacing go standard building means.

### Installing the CLI
To install the CLI you can run the following command:

```
go install github.com/leapkit/leapkit/kit@latest
```

### Generating your project
Once the CLI has been installed it can be used to build your new leapkit application. To generate a new project you can run the following command:

```
kit new superapp
```

### Setup

Once the project is generated you can setup the project by running the following command:

```
go mod download
go run ./cmd/setup
```

This will install the necessary dependencies and setup the project for development.

### Running the application

Running the application is as simple as running the following command:

```
kit dev
```

Once the application is running it should be [accessible locally](http://localhost:3000).
