package form

import (
	"net/http"
	"net/url"

	"github.com/leapkit/core/form/validate"
)

// validator is an interface that wraps the Validate method
// for validating form values.
type validator interface {
	Validate(req url.Values) validate.Errors
}

func Validate(req *http.Request, rules validator) validate.Errors {
	// Validate validates the form values with the passed rules.

	return rules.Validate(req.Form)
}
