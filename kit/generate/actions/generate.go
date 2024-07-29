package action

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	actionsFolder = "internal"
	templateFile  = `package {{.Package}}

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Name(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	if err := rw.Render("{{.Folder}}/{{.FileName}}.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}`
)

func ActionsFolder() string {
	return actionsFolder
}

func Template() string {
	return templateFile
}

func generate(f instace) error {
	// Create the folder
	if err := os.MkdirAll(filepath.Join(actionsFolder, f.folder), 0755); err != nil {
		return fmt.Errorf("error creating folder: %w", err)
	}

	// Create the handler file
	if err := f.create(".go"); err != nil {
		return err
	}

	// Create the HTML file
	if err := f.create(".html"); err != nil {
		return err
	}

	fmt.Println("Action files created successfully✅")
	return nil
}
