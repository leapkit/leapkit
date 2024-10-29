package errors

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
)

func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request, _ error) {
	renderErrorPage(w, r,
		http.StatusNotFound,
		"Something is wrong",
		"The page you are looking for was moved, removed, renamed or might never existed!",
	)
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	renderErrorPage(w, r,
		http.StatusInternalServerError,
		"Sorry, unexpected error",
		"We are working on fixing the problem, Be back soon",
	)
}

func renderErrorPage(w http.ResponseWriter, r *http.Request, statusCode int, title, description string) {
	rw := render.FromCtx(r.Context())
	rw.Set("statusCode", statusCode)
	rw.Set("title", title)
	rw.Set("description", description)
	if err := rw.Render("errors/error_page.html"); err != nil {
		server.Error(w, err, http.StatusInternalServerError)
	}
}
