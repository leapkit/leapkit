package session

import (
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
	"github.com/leapkit/leapkit/core/server/internal/writer"
)

// saver takes care of automatically saving the session
// when the response is written, this avoids having to
// call session.Save() in every handler.
type saver struct {
	*writer.ResponseWriter
	req   *http.Request
	store *sessions.Session
	moot  sync.Mutex
}

func (s *saver) Header() http.Header {
	s.saveSession()
	return s.ResponseWriter.Header()
}

func (s *saver) WriteHeader(code int) {
	s.saveSession()
	s.ResponseWriter.WriteHeader(code)
}

func (s *saver) Write(b []byte) (int, error) {
	s.saveSession()
	return s.ResponseWriter.Write(b)
}

func (s *saver) saveSession() {
	s.moot.Lock()
	defer s.moot.Unlock()

	s.store.Save(s.req, s.ResponseWriter)
}
