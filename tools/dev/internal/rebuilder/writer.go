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

// Write is the implementation of io.Writer interface for leapkit
// It writes the logs to the console with the following format:
//
// 2006-01-02 15:04:05 prefix | string(p)
func (cw *customWriter) Write(p []byte) (int, error) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	timestamp := time.Now().Format(time.DateTime)
	trailingSpaces := strings.Repeat(" ", max(maxServiceNameLen-len(cw.prefix), 0))

	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := fmt.Sprintf("%s%s %s%s|%s %s\n",
			cw.color,
			timestamp,
			cw.prefix,
			trailingSpaces,
			endColor,
			strings.TrimSuffix(scanner.Text(), "\n"),
		)

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
