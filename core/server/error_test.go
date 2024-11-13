package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leapkit/leapkit/core/server"
)

func TestErrorf(t *testing.T) {
	h := func(w http.ResponseWriter, r *http.Request) {
		server.Errorf(w, http.StatusInternalServerError, "an error occurred in %v", 10)
	}

	s := server.New()
	s.HandleFunc("/error", h)

	// Test the error response
	req, _ := http.NewRequest("GET", "/error", nil)
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d; got %d", http.StatusInternalServerError, rec.Code)
	}
}
