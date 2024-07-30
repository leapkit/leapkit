package server

import (
	"log/slog"
	"net/http"
)

// Error logs the error and sends an internal server error response.
func Error(w http.ResponseWriter, err error, HTTPStatus int) {
	slog.Error(err.Error())

	http.Error(w, err.Error(), HTTPStatus)
}
