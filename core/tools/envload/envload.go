// envload package loads .env files into the environment. To do it
// it uses an init function that reads the .env file and sets the
// variables in the environment.
package envload

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// init reads the .env file and sets the variables in the environment
// and passes it to the parseVars function, then it set the variables
// in the environment.
func init() {
	// open .env file
	file, err := os.Open(".env")
	if err != nil {
		return
	}

	for key, value := range parseVars(file) {
		err := os.Setenv(key, value)
		if err != nil {
			continue
		}
	}
}

// parseVars reads the variables from the reader and sets them
// in the environment.
func parseVars(r io.Reader) map[string]string {
	vars := make(map[string]string)
	scanner := bufio.NewScanner(r)
	var key, value string
	var isMultiLine bool

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case !isMultiLine:
			key, value, isMultiLine = parseLine(line)
		default:
			value, isMultiLine = continueMultiLineValue(value, line)
		}

		vars[key] = value
	}

	return vars
}

// parseLine parses a line from the .env file and returns the key and value
// and if the value is a multi-line value
// It returns the key, value and a boolean indicating if the value is a multi-line value
func parseLine(line string) (key, value string, isMultiLine bool) {
	pair := strings.SplitN(line, "=", 2)
	if len(pair) != 2 {
		return
	}

	key = strings.TrimSpace(pair[0])
	value = strings.TrimSpace(pair[1])

	// Check if the value is a multi-line value by checking if it starts with a quote but doesn't end with one
	if strings.HasPrefix(value, "\"") && !strings.HasSuffix(value, "\"") {
		isMultiLine = true
		value = value[1:] + "\n"
		return
	}

	// Check if the value is a multi-line value by checking if it starts and ends with a quote on the same line
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
	}

	return key, value, isMultiLine
}

// continueMultiLineValue continues a multi-line value
// by appending the line to the value until the line ends with a quote
func continueMultiLineValue(value, line string) (string, bool) {
	value += line + "\n"
	if strings.HasSuffix(line, "\"") {
		return value[:len(value)-2], false
	}

	return value, true

}
