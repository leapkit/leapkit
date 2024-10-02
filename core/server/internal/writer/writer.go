package writer

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

// ResponseWriter is a custom wrapper around http.ResponseWriter used in the server package.
// It also captures the HTTP status code and and implements the http.Flusher and http.Hijacker interfaces.
type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

// WriteHeader sets the status code and calls the WriteHeader() method of http.ResponseWriter.
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Flush method is the http.Flusher implementation of this wrapper.
// The Flush() method will be called if the wrapped http.ResponseWriter supports flushing.
func (w *ResponseWriter) Flush() {
	f, ok := w.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	f.Flush()
}

// Hijack is the implementation of the http.Hijacker trying to parse
// the wrapped http.ResponseWriter into a http.Hijacker interface.
//
// It returns an error if hijacking is not supported.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("Hijack not supported")
	}

	return h.Hijack()
}
