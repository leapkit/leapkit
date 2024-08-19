package generate

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "embed"
)

var (
	// actionsFolder is the folder where the actions are stored
	actionsFolder = "internal"

	//go:embed action.go.tmpl
	actionTemplate string
)

// Action generates a new action
func Action(name string) error {
	err := Handler(name)
	if err != nil {
		return err
	}

	// Create action.html
	folder := path.Dir(name)
	fileName := strings.ToLower(path.Base(name))
	_, err = os.Create(filepath.Join(actionsFolder, folder, fileName+".html"))
	if err != nil {
		return err
	}

	fmt.Println("Action files created successfully✅")

	return nil
}
