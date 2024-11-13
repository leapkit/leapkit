---
index: 8
title: "Session"
---

Leapkit provides a way to set a session in your app to persist values needed throughout the app. Leapkit uses the [gorilla/sessions](https://github.com/gorilla/sessions) package for this.

## Setup

To set the session in your app, you have to use the `server.WithSession` server option and pass a session secret key and the session name.

```go
s := server.New(
   server.WithSession("secret_key", "session_name"),
)
```

## Handling session values and flashes

To use the session struct within your handler, retrieve it from the context using the `session.FromCtx()` function. Then, you can manage your session values according to the `gorilla/session` package [docs](https://pkg.go.dev/github.com/gorilla/sessions). For instance:


```go
func Handler(w http.ResponseWriter, r *http.Request) {
    ss := session.FromCtx(r.Context())

    // setting a session value
    ss.Values["sessionVal"] = "value"

    // getting a session value
    myVal := ss.Values["sessionVal"].(string)


    // setting a flash message
    ss.AddFlash("welcome!")
    ss.AddFlash("peter", "username_flash")
    // {
    // 	"_flash": ["welcome!"],
    // 	"username_flash": ["peter"],
    // }

    // getting flash messages
    flashes := ss.Flashes("username_flash")
    // ["peter"]
    // ...
}
```

You can omit the `session.Save()` method **only** if you use `http.ResponseWriter` methods because the response writer is replaced by a Leapkit session implementation, which saves the current session. Otherwise, you have to use it.