---
title: "LeapKit 🚀🎒"
index: 1
---

LeapKit is a collection of packages to help you build your next application with Go. Its main purpose is to make web development as fast and simple as possible and is targeted at web developers trying to take the leap and launch their next project.

There are three components of LeapKit:
- The LeapKit Core
- The LeapKit Template
- The `kit` CLI

## LeapKit Core
The LeapKit Core module contains the Go libraries with common web features. Things such as form binding, form validating, database migrations, and routing helpers, and hot code reloading live in this package.

## LeapKit Template
The **LeapKit Template** contains a starting point folder structure using the LeapKit core that uses the LeapKit core and the standard library to provide common functionallity.

## Kit CLI
The kit CLI is a CLI tool with the development operations needed to acelerate the development of LeapKit based applications. With the `kit` CLI you can do things such as initializing a LeapKit application, running migrations, code generation and runnign the app in development mode.

The template also has some CLI commands to facilitate the web development of your apps. The template is the starting point in your development journey with LeapKit, and throught he help of the [gonew](https://github.com/golang/tools/tree/master/cmd/gonew) command is used to compose the base structure of your app.
