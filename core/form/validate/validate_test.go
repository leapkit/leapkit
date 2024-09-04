package validate

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	verrs := Errors{
		"name":  []error{fmt.Errorf("Name is required")},
		"email": []error{fmt.Errorf("Email is required"), fmt.Errorf("Email is not valid")},
	}

	t.Run("hasAny", func(t *testing.T) {
		if !verrs.HasAny() {
			t.Fatal("expected errors to be present")
		}

		xerrs := Errors{}
		if xerrs.HasAny() {
			t.Fatal("expected no errors")
		}
	})

	t.Run("Has field", func(t *testing.T) {
		if !verrs.Has("name") {
			t.Fatal("expected name error")
		}

		if verrs.Has("bio") {
			t.Fatal("expected no bio error")
		}

		if !verrs.Has("email") {
			t.Fatal("expected email error")
		}
	})

	t.Run("ErrorStringFor", func(t *testing.T) {
		if verrs.ErrorStringFor("name") != "Name is required" {
			t.Fatal("expected name error")
		}

		if verrs.ErrorStringFor("bio") != "" {
			t.Fatal("expected no bio error")
		}

		if verrs.ErrorStringFor("email") != "Email is required. Email is not valid" {
			t.Fatal("expected email error")
		}
	})

}
