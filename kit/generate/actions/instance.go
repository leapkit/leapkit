package action

import (
	"html/template"
	"os"
	"path/filepath"
)

type instace struct {
	pckg     string
	folder   string
	fileName string
}

func New(path []string) instace {
	fileName := path[len(path)-1] // The last element is the file name
	folder := path[:len(path)-1]  // The folder is the path without the file name
	pckg := folder[len(folder)-1] // The package is the last element of the folder

	return instace{
		pckg:     pckg,
		folder:   filepath.Join(folder...),
		fileName: fileName,
	}
}

func (f instace) Generate() error {
	return generate(f)
}

func (f instace) create(ext string) error {
	file, err := os.Create(filepath.Join(ActionsFolder(), f.folder, f.fileName+ext))
	if err != nil {
		return err
	}

	defer file.Close()

	if ext == ".html" {
		return nil
	}

	template := template.Must(template.New("handler").Parse(Template()))
	err = template.Execute(file, map[string]string{
		"Package":  f.pckg,
		"FileName": f.fileName,
		"Folder":   f.folder,
	})

	if err != nil {
		return err
	}

	return nil

}
