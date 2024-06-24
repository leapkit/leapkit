package assets

import (
	"io/fs"
	"os"
	"sync"
)

type manager struct {
	embedded fs.FS
	folder   fs.FS

	outputFolder string
	inputFolder  string

	servingPath string

	fmut       sync.Mutex
	fileToHash map[string]string
	HashToFile map[string]string
}

// NewManager returns a new manager that wraps the given embed.FS and the input and output folders.
func NewManager(embedded fs.FS) *manager {
	// TODO: options to change:
	// - input
	// - output.
	// - serving path.
	return &manager{
		embedded: embedded,
		folder:   os.DirFS("public"),

		inputFolder:  "internal/assets",
		outputFolder: "public",
		servingPath:  "/public/*",

		fileToHash: map[string]string{},
		HashToFile: map[string]string{},
	}
}
