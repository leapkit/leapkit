package rebuilder

import (
	"path/filepath"
	"strings"
)

var config = &configuration{}

type configuration struct {
	path string

	extensionsToWatch []string
	excludedPaths     []string

	runners []func()
}

func (c *configuration) isExcludedPath(path string) bool {
	for _, p := range c.excludedPaths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

func (c *configuration) isWatchedExtension(path string) bool {
	for _, ext := range c.extensionsToWatch {
		if filepath.Ext(path) != ext {
			continue
		}

		return true
	}

	return false
}
