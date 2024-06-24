package form_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/leapkit/core/form"
	"github.com/leapkit/core/form/validate"
)

func TestValidate(t *testing.T) {
	reqFromParams := func(params url.Values) *http.Request {
		req := httptest.NewRequest("POST", "/", bytes.NewBufferString(params.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ParseForm()

		return req
	}

	emailExp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	t.Run("Valid simple request", func(t *testing.T) {
		req := reqFromParams(url.Values{
			"name": {"John"},
		})

		rules := validate.Fields(
			validate.Field("name", validate.Required()),
		)

		errs := form.Validate(req, rules)
		if len(errs) > 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Invalid simple request", func(t *testing.T) {
		req := reqFromParams(url.Values{
			"name": {""},
		})

		rules := validate.Fields(
			validate.Field("name", validate.Required()),
		)

		errs := form.Validate(req, rules)
		if len(errs) == 0 {
			t.Fatalf("expected errors, got none")
		}
	})

	t.Run("Valid multipoe fields validation", func(t *testing.T) {
		req := reqFromParams(url.Values{
			"first_name":    {"antonio"},
			"last_name":     {"pagano"},
			"email_address": {"a@pagano.id"},
		})

		rules := validate.Fields(
			validate.Field("first_name", validate.Required()),
			validate.Field("last_name", validate.Required()),

			validate.Field(
				"email_address",
				validate.Required(),
				validate.MatchRegex(emailExp),
			),
		)

		errs := form.Validate(req, rules)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("Invalid multipoe fields validation", func(t *testing.T) {
		req := reqFromParams(url.Values{
			"first_name":    {"antonio"},
			"last_name":     {"pagano"},
			"email_address": {"a"},
		})

		rules := validate.Fields(
			validate.Field("first_name", validate.Required()),
			validate.Field("last_name", validate.Required()),

			validate.Field("email_address", validate.Required(), validate.MatchRegex(emailExp, "invalid email address")),
		)

		errs := form.Validate(req, rules)
		if len(errs) == 0 {
			t.Fatalf("expected errors, got none")
		}
	})
}
