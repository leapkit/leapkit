package server

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/leapkit/leapkit/core/server/internal/response"
)

// baseMiddleware is a list that holds the middleware list that will be executed
// at the beginning of a client request.
var baseMiddleware = []Middleware{
	setValuer,
	requestID,
	logger,
}

// Middleware is a function that receives a http.Handler and returns a http.Handler
// that can be used to wrap the original handler with some functionality.
type Middleware func(http.Handler) http.Handler

func requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "requestID", time.Now().UnixNano()))
		next.ServeHTTP(w, r)
	})
}

// logger is a middleware that logs the request method and URL
// and the time it took to process the request.
func logger(next http.Handler) http.Handler {
	logger := slog.Default()
	if os.Getenv("GO_ENV") == "production" {
		// Using json logger in production
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw, ok := w.(*response.Writer)
		if !ok {
			lw = &response.Writer{ResponseWriter: w}
		}

		defer func() {
			status := cmp.Or(lw.Status, http.StatusOK)
			logLevel := slog.LevelInfo

			if status >= http.StatusInternalServerError {
				logLevel = slog.LevelError
			}

			logger.Log(r.Context(), logLevel, "", "method", r.Method, "status", status, "url", r.URL.Path, "took", time.Since(start))
		}()

		next.ServeHTTP(lw, r)
	})
}

// recoverer is a middleware that recovers from panics and logs the error.
// The error stack trace is printed only when the application is in 'development' mode.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lw, ok := w.(*response.Writer)
		if !ok {
			lw = &response.Writer{ResponseWriter: w}
		}

		defer func() {
			if err := recover(); err != nil || lw.Status >= http.StatusInternalServerError {
				slog.Error("panic", "error", err, "method", r.Method, "url", r.URL.Path)

				if cmp.Or(os.Getenv("GO_ENV"), "development") == "development" {
					os.Stderr.WriteString(fmt.Sprint(err, "\n"))
					debug.PrintStack()
				}

				if lw.Status == 0 {
					w.WriteHeader(http.StatusInternalServerError)
				}

				errorHandlerMap[http.StatusInternalServerError](lw, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(lw, r)
	})
}
