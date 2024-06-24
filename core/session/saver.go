package session

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/sessions"
)

// saver takes care of automatically saving the session
// when the response is written, this avoids having to
// call session.Save() in every handler.
type saver struct {
	w http.ResponseWriter

	req   *http.Request
	store *sessions.Session
	moot  sync.Mutex
}

func (s *saver) Header() http.Header {
	s.moot.Lock()
	defer s.moot.Unlock()

	s.store.Save(s.req, s.w)
	return s.w.Header()
}

func (s *saver) WriteHeader(code int) {
	s.moot.Lock()
	defer s.moot.Unlock()

	s.store.Save(s.req, s.w)
	s.w.WriteHeader(code)
}

func (s *saver) Write(b []byte) (int, error) {
	s.moot.Lock()
	defer s.moot.Unlock()

	s.store.Save(s.req, s.w)
	n, err := s.w.Write(b)
	return n, err
}

func (s *saver) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := s.w.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}

	return h.Hijack()
}
