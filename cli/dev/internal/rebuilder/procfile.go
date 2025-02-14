package rebuilder

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type entry struct {
	ID      int
	Name    string
	Command string
}

func readProcfile(path string) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var entries []entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}

		exists := slices.ContainsFunc(entries, func(e entry) bool {
			return e.Name == parts[0]
		})

		if exists {
			continue
		}

		entries = append(entries, entry{
			ID:      len(entries),
			Name:    strings.TrimSpace(parts[0]),
			Command: strings.TrimSpace(parts[1]),
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Procfile: %w", err)
	}

	return entries, nil
}
