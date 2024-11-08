package errors

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed page.html
var htmlTemplate string

func NotFoundPage() string {
	return build(
		http.StatusNotFound,
		"Something went wrong",
		"The page you are looking for was moved, removed, renamed or might never existed!",
	)
}

func InternalServerErrorPage() string {
	return build(
		http.StatusInternalServerError,
		"We're fixing it",
		"This page is having some technical hiccups. We know about the problem and we're working to get this back to normal quickly",
	)
}

func build(status int, title, description string) string {
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
