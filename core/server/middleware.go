package server

import (
	"cmp"
	"context"
	"log/slog"
	"net/http"
	"time"
)

// baseMiddleware is a list that holds the middleware list that will be executed
// at the beginning of a client request.
var baseMiddleware = []Middleware{
	setValuer,
	requestID,
	recoverer,
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

// loggerWriter is a wrapper around http.ResponseWriter that keeps track of the status code
type loggerWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader methods overrides the ResponseWriter.WriteHeader method to capture the status code.
func (lw *loggerWriter) WriteHeader(statusCode int) {
	lw.status = statusCode
	lw.ResponseWriter.WriteHeader(statusCode)
}

// logger is a middleware that logs the request method and URL
// and the time it took to process the request.
func logger(next http.Handler) http.Handler {
	logger := slog.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := &loggerWriter{ResponseWriter: w}

		defer func() {
			lw.status = cmp.Or(lw.status, http.StatusOK)
			logLevel := slog.LevelInfo

			if lw.status >= http.StatusInternalServerError {
				logLevel = slog.LevelError
			}

			logger.Log(r.Context(), logLevel, "", "method", r.Method, "status", lw.status, "url", r.URL.Path, "took", time.Since(start))
		}()

		next.ServeHTTP(lw, r)
	})
}

// recoverer is a middleware that recovers from panics and logs the error.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err, "method", r.Method, "url", r.URL.Path)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
