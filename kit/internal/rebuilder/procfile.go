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

func readProcfile() ([]entry, error) {
	f, err := os.Open("Procfile")
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
			ID:      len(entries) + 1,
			Name:    parts[0],
			Command: parts[1],
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading Procfile: %w", err)
	}

	slices.SortFunc(entries, func(a, b entry) int {
		return strings.Compare(a.Name, b.Name)
	})

	var sorted []entry
	for _, e := range entries {
		maxServiceNameLen = max(maxServiceNameLen, len(e.Name))

		if e.Name == "app" {
			sorted = append([]entry{e}, sorted...)
			continue
		}

		sorted = append(sorted, e)
	}

	return sorted, nil
}
