package assets

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// manager watches the input folder and copies all files to the output folder.
// It also watches for changes in the input folder and copies the files again.
func Watch(inputFolder, outputFolder string) func() {
	return func() {
		err := copyAll(inputFolder, outputFolder)
		if err != nil {
			fmt.Println("[error]", err.Error())

			return
		}

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("[error] error creating watcher: %s\n", err.Error())

			return
		}

		// Add all folders within the assets folder to the watcher.
		err = filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
			return watcher.Add(path)
		})

		if err != nil {
			fmt.Printf("error adding files to watcher stopping the assets watcher: %s\n", err)

			return
		}

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						continue
					}

					needsCopy := event.Has(fsnotify.Create) || event.Has(fsnotify.Write) || event.Has(fsnotify.Rename)
					if !needsCopy {
						continue
					}

					err := copyAll(inputFolder, outputFolder)
					if err != nil {
						log.Println(err)
					}

					if event.Has(fsnotify.Create) {
						watcher.Add(event.Name)
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					log.Println("error:", err)
				}
			}
		}()

		<-make(chan struct{})
	}

}

// copyAll files from the input folder to the output folder.
func copyAll(inputFolder, outputFolder string) error {
	err := filepath.Walk(inputFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get the relative path of the file
		relativePath, err := filepath.Rel(inputFolder, path)
		if err != nil {
			return err
		}

		// Create the destination folder if it doesn't exist
		destFolder := filepath.Join(outputFolder, filepath.Dir(relativePath))
		err = os.MkdirAll(destFolder, os.ModePerm)
		if err != nil {
			return err
		}

		// Copy the file to the destination folder
		destPath := filepath.Join(destFolder, filepath.Base(relativePath))
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("error copying files: %w", err)
	}

	return nil
}
