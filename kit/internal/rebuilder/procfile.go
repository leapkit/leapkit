package rebuilder

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
)

type entry struct {
	ID      int
	Name    string
	Command string
}

func procfile() ([]entry, error) {
	f, err := os.Open("Procfile")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	rgx := regexp.MustCompile(`^([\w-]+):\s*(.+)$`)
	var entries []entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := rgx.FindStringSubmatch(scanner.Text())
		if len(parts) == 3 {
			entries = append(entries, entry{
				ID:      len(entries) + 1,
				Name:    parts[1],
				Command: parts[2],
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("[error] Error reading Procfile: %w", err)
	}

	slices.SortFunc(entries, func(a, b entry) int {
		return strings.Compare(a.Name, b.Name)
	})

	var sorted []entry
	for _, e := range entries {
		if e.Name == "app" {
			sorted = append([]entry{e}, sorted...)
			continue
		}

		sorted = append(sorted, e)
	}

	return sorted, nil
}
