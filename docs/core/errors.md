---
index: 6
title: "Errors"
---

Leapkit provides a way to customize the response message for specific HTTP error statuses like `404 not found` or `500 Internal Server Error`.

## Usage

To set custom error messages, use the `WithErrorMessage` server option. This option accepts two parameters: the HTTP status error code and the message or template to return.

```go
r := server.New(
    // Setting message for 404 status.
	server.WithErrorMessage(
        http.StatusNotFound,
        "Oops! We couldn't find the page you were looking for",
    ),
    // Setting message for 500 status.
	server.WithErrorMessage(
        http.StatusInternalServerError,
        "There were some technical issues while processing your request",
    ),
)
```

## HTML error pages

Leapkit also supports returning HTML error pages. You can pass HTML templates as strings, and they will be rendered directly in the browser. Leapkit automatically sets the `Content-Type` response header based on the content passed to `WithErrorMessage`:

- If the message is plain text, the `Content-Type` is set to `text/plain`.
- If the message is HTML, the `Content-Type` is set to `text/html`.


```go
r := server.New(
    server.WithErrorMessage(http.StatusNotFound, notFoundPage())
)

func notFoundPage() string {
    return `
    <html>
        <head>
            <title>404 Not Found</title>
        </head>
        <body>
            <h1>Oops! We couldn't find the page you were looking for</h1
        </body>
    </html>`
}
```

## Error handling

When an error needs to be handled in handlers you usually need to return the error message to the client and return. This is done typically the following way:

```go
func (w http.ResponseWriter, r *http.Request) {
...
if err != nil {
	http.Error(w, err.Error(), http.StatusNotFound)
	slog.Error("Error happened", "error", err)

	return
}

...
}
```

### server.Error

Leapkit provides a helper function to return an error message with a specific status code. This function is useful and acts as a shortcut for the `http.Error` function and logging the error.

```go
server.Error(w, err, http.StatusNotFound)
```

### server.Errorf

In some cases you would want to format the error message before returning it. The `server.Errorf` function is a helper function that formats the error message and returns it with a specific status code.

```go
server.Errorf(w, http.StatusNotFound, "Error happened: %v", err.Error())
```
