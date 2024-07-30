package action

import (
	"html/template"
	"os"
	"path/filepath"
)

type instance struct {
	pckg     string
	folder   string
	fileName string
}

func New(path []string) instance {
	var pckg string = ActionsFolder()
	fileName := path[len(path)-1] // The last element is the file name

	if len(path) == 1 {
		return instance{
			pckg:     pckg,
			folder:   "",
			fileName: fileName,
		}
	}

	folderSlice := path[:len(path)-1]             // The folder is the path without the file name
	pckg = folderSlice[len(folderSlice)-1]        // The package is the last element of the folder
	folder := filepath.Join(folderSlice...) + "/" // Join folder elements with "/" and append "/"
	return instance{
		pckg:     pckg,
		folder:   folder,
		fileName: fileName,
	}
}

func (f instance) Generate() error {
	return generate(f)
}

func (f instance) create(ext string) error {
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
