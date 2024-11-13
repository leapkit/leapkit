---
title: Rendering
index: 7
---

Leapkit offers the `render` package which handles and renders the HTML files in your project.

## Setting up

Philosophy of the render engine is that it lives in the application context and is called within each handler when you need to render an HTML page. So, to setup this in your application you need to set the engine in context by using its `Middleware` function.


### The `Middleware` function

> **NOTE:** You have the option to use the `InCtx` function as an alternative. This one perform the same action as `Middleware`.

This is the way that the `render` package provides to initialize the render engine in the context of your application, this function receives the file system (`fs.FS`) which is pointing to your HTML files dir. This function returns a middleware that must be used by the router.

```go
import "github.com/leapkit/leapkit/core/render"

var (

    // Alternatively...
    // renderMW = render.InCtx(templates.FS)
    renderMW = render.Middleware(templates.FS)
)
// ...
func appRoutes() {
    s := server.New()

    s.Use(renderMW)
    //...
}
```

### Options

The `Middleware` function allows some options that you can set when this is being called:

#### WithDefaultLayout

This option allows you to set the base template on which all HTML files will be built.

```text
├── your/path/to/templates/
│	├── default
│	│   └── application.html
│	└── templates.go  // FS
```

```go
// renderMW = render.InCtx(...
renderMW = render.Middleware(templates.FS,
    render.WithDefaultLayout("default/application.html"),
)
```

#### WithHelpers

This is another option which allows you to set helpers that you can use in your templates.

```go
myHelpers = map[string]any{
    //...
}

// renderMW = render.InCtx(...
renderMW = render.Middleware(templates.FS,
    render.WithHelpers(myHelpers),
)
```

## Getting the render engine

To get the render engine from context, you can use the `FromCtx()` function which receives a context parameter.

```go
func appHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())
}
```

## Setting variables in context

You can set variables that will be used within HTML pages by using the `Set` render engine method

```go
func appHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())

    rw.Set("myVariable", value)
}

```


## Rendering HTML pages

The render engine counts with three different ways to render an HTML content, let's take a look by using this example:


```text
├── templates/
│	├── layouts/
│	│   ├── default.html
│	│   └── simple.html
│	├── users/
│	│   └── list.html
│	│
│	└── templates.go
│
├── internal/
│	├── users/
│	│   └── users.go
│	│
│	└── routes.go
```

```html
<!-- templates/layouts/default.html -->
<html>
    <body class="h-full">
        <header class="header-class"></header>
        <%= yield %>
        <footer class="footer-class"></footer>
    </body>
</html>
```

```html
<!-- templates/layouts/simple.html -->
<html>
    <body>
        <%= yield %>
    </body>
</html>
```

```html
<!-- templates/users/list.html -->
<ul>
    <%= for (user) in users { %>
        <li><%= user %></li>
    <% } %>
</ul>
```

```go
// templates/templates.go
package templates

import "embed"

//go:embed */*.html
var FS embed.FS
```

```go
// templates/users/users.go
package users


var usersList = []string{
    "John Smith",
    "Smart Josh",
}

func UserListHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())

    rw.Set("users", usersList)
    // ...
}
```

```go
// internal/routes.go
package internal

import (
    "app/templates"
    "app/internal/users"
    "github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
)

var (
    // renderMW = render.InCtx(...
    renderMW = render.Middleware(templates.FS,
        render.WithDefaultLayout("layouts/default.html"),
    )
)

func Mount(r server.Router) {
    r.Use(renderMW)
    r.HandleFunc("GET /users/{$}", users.UserListHandler)
}
```

### Render

This method renders the HTML page on the base layout.

```go
func UserListHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())

    rw.Set("users", usersList)
    if err := rw.Render("users/index.html"); err != nil {
        // handle the err
    }
    // ...
}
```

```html
<!-- Result -->
<html>
    <body class="h-full">
        <header class="header-class"></header>
        <ul>
            <li>John Smith</li>
            <li>Smart Josh</li>
        </ul>
        <footer class="footer-class"></footer>
    </body>
</html>
```

### RenderWithLayout

This method renders the HTML page on an specific layout.

```go
func UserListHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())

    rw.Set("users", usersList)
    if err := rw.RenderWithLayout("users/index.html", "layouts/simple.html"); err != nil {
        // handle the err
    }
    // ...
}
```

```html
<!-- Result -->
<html>
    <body>
        <ul>
            <li>John Smith</li>
            <li>Smart Josh</li>
        </ul>
    </body>
</html>
```

### RenderClean

The following method renders the HTML file without wrapping it within any layout.

```go
func UserListHandler(w http.ResponseWriter, r *http.Request) {
    rw := render.FromCtx(r.Context())

    rw.Set("users", usersList)
    if err := rw.RenderClean("users/index.html"); err != nil {
        // handle the err
    }
    // ...
}
```

```html
<!-- Result -->
<ul>
    <li>John Smith</li>
    <li>Smart Josh</li>
</ul>
```
