package assets

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (m *manager) HandlerPattern() string {
	return m.servingPath
}

func (m *manager) HandlerFn(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, m, strings.TrimPrefix(r.URL.Path, m.handlerPrefix()))
}

func (m *manager) Open(name string) (file fs.File, err error) {
	ext := filepath.Ext(name)
	if ext == ".go" {
		return nil, os.ErrNotExist
	}

	// Converting hashed into original file name
	smp := m.HashToFile[name]
	if smp != "" {
		name = smp
	}

	fn := m.embedded.Open
	if env := os.Getenv("GO_ENV"); env == "development" {
		fn = m.folder.Open
	}

	name = strings.TrimPrefix(name, m.handlerPrefix())

	file, err = fn(name)
	return file, err
}

func (m *manager) ReadFile(name string) ([]byte, error) {
	x, err := m.Open(name)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(x)
}

func (m *manager) handlerPrefix() string {
	return strings.TrimSuffix(m.servingPath, "*")
}
