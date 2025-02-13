package rebuilder

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

var maxServiceNameLen = 4

var colors = []string{
	"\033[92m", // Green
	"\033[93m", // Yellow
	"\033[94m", // Blue
	"\033[95m", // Magenta
	"\033[96m", // Cyan
}

const endColor = "\033[0m"

type customWriter struct {
	mu     sync.Mutex
	writer io.Writer
	prefix string
	color  string
}

func (cw *customWriter) Write(p []byte) (int, error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	timestamp := time.Now().Format(time.DateTime)
	trailingSpaces := strings.Repeat(" ", maxServiceNameLen-len(cw.prefix))

	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}

		line := fmt.Sprintf("%s%s %s%s |%s", cw.color, timestamp, cw.prefix, trailingSpaces, endColor)
		line += " " + strings.TrimSuffix(scanner.Text(), "\n") + "\n"

		if _, err := cw.writer.Write([]byte(line)); err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

func wrap(writer io.Writer, e entry) io.Writer {
	return &customWriter{
		writer: writer,
		prefix: e.Name,
		color:  colors[e.ID%len(colors)],
	}
}
