package generate

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed handler.go.tmpl
	handlerTemplate string
)

func Handler(name string) error {
	path := strings.Split(name, string(filepath.Separator))
	actionPackage := "internal"
	fileName := path[len(path)-1] // file name is the last part of the path
	if len(path) > 1 {
		actionPackage = path[len(path)-2] // package name is the second to last part of the path
	}

	folder := strings.Join(path[:len(path)-1], string(filepath.Separator)) // folder is everything but the last part of the path

	// Create the folder
	if actionPackage != "internal" {
		if err := os.MkdirAll(filepath.Join(actionsFolder, folder), 0755); err != nil {
			return fmt.Errorf("error creating folder: %w", err)
		}
	}

	// Create action.go
	file, err := os.Create(filepath.Join(actionsFolder, folder, fileName+".go"))
	if err != nil {
		return err
	}

	defer file.Close()
	template := template.Must(template.New("handler").Parse(handlerTemplate))
	fileName = cases.Title(language.English).String(filepath.Base(fileName))
	err = template.Execute(file, map[string]string{
		"Package":  actionPackage,
		"FuncName": fileName,
	})

	if err != nil {
		return err
	}

	fmt.Println("Handler file created successfully✅")
	return nil
}
