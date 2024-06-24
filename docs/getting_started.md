---
title: Getting Started
index: 2
---

To get started with LeapKit the [template](https://github.com/leapkit/template) repository provides a basic template for building web applications with Go, HTMX and Tailwind CSS. The template generates useful commands to help you develop, build and deploy your application.

### Generating your project

To generate your project you can use `gonew` to copy the template and generate the necessary files for your project. The following command will generate a new project called `superapp` in the current directory.

```
go run rsc.io/tmp/gonew@latest github.com/leapkit/template@v1.1.8 superapp
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
go run ./cmd/dev
```

Once the application is running it should be [accessible locally](http://localhost:3000).
