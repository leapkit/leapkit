package server

import (
	"cmp"
	"log/slog"
	"net/http"
)

// Error writes an HTTP error response, logging the error message.
// Unlike http.Error, this function determines the Content-Type dynamically
// depending on whether a message is found in the errorMessageMap for the given HTTPStatus.
// If no error message is registered, it defaults to the error's message content type.
func Error(w http.ResponseWriter, err error, HTTPStatus int) {
	slog.Error(err.Error())

	content := []byte(cmp.Or(errorMessageMap[HTTPStatus], err.Error()))

	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", http.DetectContentType(content))
	h.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(HTTPStatus)
	w.Write(content)
}

// errorMessageMap holds custom error messages based on HTTP status codes.
var errorMessageMap = map[int]string{}
