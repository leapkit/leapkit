package generate

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed handler.go.tmpl
	handlerTemplate string
)

func handlerName(name string) string {
	var g []string
	p := strings.Fields(name)
	for _, value := range p {
		g = append(g, cases.Title(language.English, cases.NoLower).String(value))
	}

	return strings.Join(g, "")
}

func Handler(name string) error {
	actionPackage := "internal"

	folder := path.Dir(name)
	if len(folder) > 0 && folder != "." {
		parts := filepath.SplitList(folder)
		actionPackage = parts[len(parts)-1]
		actionPackage = cases.Lower(language.English).String(folder)
	}

	// Create the folder
	if err := os.MkdirAll(filepath.Join(actionsFolder, folder), 0755); err != nil {
		return fmt.Errorf("error creating folder: %w", err)
	}

	// Create action.go
	fileName := strings.ToLower(path.Base(name))
	file, err := os.Create(filepath.Join(actionsFolder, folder, fileName+".go"))
	if err != nil {
		return err
	}

	defer file.Close()
	template := template.Must(template.New("handler").Parse(handlerTemplate))
	err = template.Execute(file, map[string]string{
		"Package":  actionPackage,
		"FuncName": handlerName(name),
	})

	if err != nil {
		return err
	}

	fmt.Println("Handler file created successfully✅")
	return nil
}
