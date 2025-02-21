package assets

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strings"
)

// PathFor returns the fingerprinted path for a given
// file. If the path passed contains the hash it will
// return the same path.

// filename to open should be the file without the prefix
// filename for the map should be the file without the prefix
// filename returned should be the file with the prefix
func (m *manager) PathFor(name string) (string, error) {
	normalized := m.normalize(name)
	if hashed, ok := m.fileToHash[normalized]; ok {
		return path.Join(m.servingPath, hashed), nil
	}

	// Compute the hash of the file
	bb, err := m.ReadFile(normalized)
	if err != nil {
		return "", fmt.Errorf("could not open %q: %w", normalized, os.ErrNotExist)
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

	return path.Join("/", m.servingPath, filename), nil
}

func (m *manager) normalize(name string) string {
	name = strings.TrimPrefix(path.Clean(name), "/")
	servingPath := strings.TrimPrefix(path.Clean(m.servingPath), "/")

	name = strings.TrimPrefix(path.Clean(name), servingPath)

	return strings.TrimPrefix(name, "/")
}
