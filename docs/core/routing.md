---
index: 1
title: "Routing"
---

Leapkit uses Go 1.22 routing as the core of its routing system. This means that you can use the same routing system as you would in a Go application. This is done by using the `http.HandlerFunc` type as the handler for the routes.

```go
r.HandleFunc("GET /{$}", channels.Default)
r.Group("/settings/", func(r server.Router) {
	r.Use(users.OnlyAdmin)

	r.HandleFunc("GET /edit", settings.Edit)
	r.HandleFunc("GET /join-token/generate", settings.GenerateJoinToken)
	r.HandleFunc("PUT /update", settings.Update)
})
```

## LeapKit Server
The leapkit server is a wrapper around the Go `http.Server` struct. It provides some extra abilities to the server, such as the ability to group routes and middleware.

The `server` paackage is where all these features live. You can create a new server by calling the `server.New` function.

```go
package main

import (
	"github.com/leapkit/leapkit/server"
)

func main() {
	s := server.New()

	fmt.Println("Server started at", s.Addr())
	if err := http.ListenAndServe(s.Addr(), s.Handler()); err != nil {
		fmt.Println(err)
	}
}
```

Returned Router instance configured with a default router so you can add handlers just like you would in a Go application.

### Built in middleware

The server has some built-in middleware that you can use to add some extra functionality to your server.

- Logging
- Panic recovering
- RequestID
- ValueSetter **

## Router options
The router returned by the `server.New` function can receive some options that you can use to configure the server.

```go
// Initializing the server with some options
s := server.New(
	server.WithHost("localhost"),    // Set the host of the server
	server.WithPort("8080"),         // Set the port of the server
)
```

### WithHost
WithHost allows to configure the host of the server. By default its `0.0.0.0`.

### WithPort
WithPort allows to specify the port of the server. By default its `3000`.

### WithSession
WithSession allows to set a new session into the server. [Read more](/core/session.html).

### WithAssets
WithAssets allows to set assets into the server. [Read more](/core/assets.html).

### WithErrorMessage
WithErrorMessage allows you to set your custom 404 or 500 messages. [Read more](/core/errors.html).

## Middleware
The Router returned by the `server.New` function has a `Use` method that allows you to add middleware to the server.

```go
// headerMW is a middleware that adds a header to the response
func headerMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Hello", "World")

		next.ServeHTTP(w, r)
	})
}

// ...
	s.Use(headerMW)
	s.HandleFunc("/hello", helloHandler)

	fmt.Println("Server started at", s.Addr())
// ...
```

## Grouping Routes
The Router returned by the `server.New` function has a `Group` method that allows you to group routes together, this is useful to have a better organization of your routes.

```go
package main

import (
	"net/http"

	"github.com/leapkit/leapkit/server"
)

func main() {
	s := server.New()

	// The API group of routes
	s.Group("/api", func(r *server.Router) {
		r.HandleFunc("GET /hello", helloHandler)
		r.HandleFunc("POST /hello", createHelloHandler)
		r.HandleFunc("DELETE /hello", deleteHelloHandler)
	})

	fmt.Println("Server started at", s.Addr())
	if err := http.ListenAndServe(s.Addr(), s.Handler()); err != nil {
		fmt.Println(err)
	}
}
```

## Folder Serving

The Router returned by the `server.New` function has a `ServeFiles` method that allows you to serve files from a folder or any other io.FS.

```go
package main

import (
	"net/http"

	"github.com/leapkit/leapkit/server"
)

func main() {
	s := server.New()

	// ... Other routes

	// Serve files from the public folder
	s.ServeFiles("/public", http.Dir("public"))

	fmt.Println("Server started at", s.Addr())
	if err := http.ListenAndServe(s.Addr(), s.Handler()); err != nil {
		fmt.Println(err)
	}
}
```
