package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"
)

var normalized = func(name string) string {
	name = strings.TrimPrefix(name, "/public/")
	name = strings.TrimPrefix(name, "public/")

	return name
}

var withPrefix = func(name string) string {
	return path.Join("/public/", name)
}

// PathFor returns the fingerprinted path for a given
// file. If the path passed contains the hash it will
// return the same path.

// filename to open should be the file without the prefix
// filename for the map should be the file without the prefix
// filename returned should be the file with the prefix
func (m *manager) PathFor(fname string) (string, error) {
	normalized := normalized(fname)
	result := m.fileToHash[normalized]
	if result != "" {
		return withPrefix(result), nil
	}

	// Compute the hash of the file
	bb, err := m.ReadFile(normalized)
	if err != nil {
		return "", fmt.Errorf("could not open %s: %w", normalized, os.ErrNotExist)
	}

	hash := md5.Sum(bb)
	hashString := hex.EncodeToString(hash[:])

	// Add the hash to the filename
	ext := path.Ext(normalized)
	filename := strings.TrimSuffix(normalized, ext)
	filename += "-" + hashString + ext

	m.fmut.Lock()
	defer m.fmut.Unlock()
	m.fileToHash[normalized] = filename
	m.HashToFile[filename] = normalized

	return withPrefix(filename), nil
}
