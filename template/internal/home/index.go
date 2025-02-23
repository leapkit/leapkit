package home

import (
	"net/http"

	"go.leapkit.dev/core/render"
	"go.leapkit.dev/core/server"
)

// Renders the home page of the application.
func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	err := rw.Render("home/index.html")
	if err != nil {
		server.Errorf(w, http.StatusInternalServerError, "error rendering template: %s", err.Error())
	}
}
