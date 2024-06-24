package envload

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseVars(t *testing.T) {

	t.Run("simple one", func(t *testing.T) {
		r := strings.NewReader("KEY=value\n")
		vars := parseVars(r)

		if vars["KEY"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY"])
		}
	})

	t.Run("multiple", func(t *testing.T) {
		vars := parseVars(strings.NewReader(`
			KEY=value
			KEY2=value
		`))

		if vars["KEY"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY"])
		}

		if vars["KEY2"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY"])
		}
	})

	t.Run("quotes", func(t *testing.T) {
		vars := parseVars(strings.NewReader(`
			KEY="value"
		`))

		if vars["KEY"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY"])
		}
	})

	t.Run("multiple equals sign", func(t *testing.T) {
		vars := parseVars(strings.NewReader(`
			KEY="value=with=equals"
		`))

		if vars["KEY"] != "value=with=equals" {
			t.Errorf("Expected value to be 'value=with=equals', got %s", vars["KEY"])
		}
	})

	t.Run("multiple lines", func(t *testing.T) {
		value := `-----BEGIN RSA PRIVATE KEY-----
		MIIEpAIBAAKCAQEAqTmwQppL07nBl/0TEQ5sHcqj/Iz9BmuaaEu26jMXYt1QttHn
		-----END RSA PRIVATE KEY-----`

		vars := parseVars(strings.NewReader(fmt.Sprintf(`
		KEY="%s"`, value)))

		if vars["KEY"] != value {
			t.Errorf("Expected value to be %s', got %s", value, vars["KEY"])
		}
	})

	t.Run("comments", func(t *testing.T) {
		vars := parseVars(strings.NewReader(`
			# This is a comment
			KEY=value
		`))

		if vars["KEY"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY"])
		}
	})

	t.Run("multiple lines with another variable", func(t *testing.T) {
		value := `-----BEGIN RSA PRIVATE KEY-----
		MIIEpAIBAAKCAQEAqTmwQppL07nBl/0TEQ5sHcqj/Iz9BmuaaEu26jMXYt1QttHn
		-----END RSA PRIVATE KEY-----`

		vars := parseVars(strings.NewReader(fmt.Sprintf(`
		KEY="%s"
		KEY2=value`, value)))

		if vars["KEY"] != value {
			t.Errorf("Expected value to be %s', got %s", value, vars["KEY"])
		}

		if vars["KEY2"] != "value" {
			t.Errorf("Expected value to be 'value', got %s", vars["KEY2"])
		}
	})

}
