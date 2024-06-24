package home

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	err := rw.Render("home/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
