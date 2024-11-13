package server

import (
	"cmp"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

var (
	//go:embed error.html
	htmlTemplate string

	// errorMessageMap holds custom error messages based on HTTP status codes.
	errorMessageMap = map[int]string{
		http.StatusNotFound: errorTemplate(
			http.StatusNotFound,
			"Something went wrong",
			"The page you are looking for was moved, removed, renamed or might never existed!",
		),
		http.StatusInternalServerError: errorTemplate(
			http.StatusInternalServerError,
			"We're fixing it",
			"This page is having some technical hiccups. We know about the problem and we're working to get this back to normal quickly",
		),
	}
)

func errorTemplate(status int, title, description string) string {
	t, err := template.New("error.html").Parse(htmlTemplate)
	if err != nil {
		return err.Error()
	}

	data := map[string]any{
		"status":      status,
		"title":       title,
		"description": description,
	}

	var builder strings.Builder
	if err := t.Execute(&builder, data); err != nil {
		return err.Error()
	}

	return builder.String()
}

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

// Errorf is a convenient function to write an HTTP error response with a formatted message
// it under the hood uses fmt.Errorf and then calls Error to write the response and log the
// error.
func Errorf(w http.ResponseWriter, HTTPStatus int, message string, args ...any) {
	Error(w, fmt.Errorf(message, args...), HTTPStatus)
}
