package validate

import (
	"net/url"
	"strings"
)

// Field validation specifies the rules for that field.
func Field(field string, rules ...ValidatorFn) fieldValidation {
	return fieldValidation{
		Field:      field,
		Validators: rules,
	}
}

// Fields is a convenience method to create a set of field validations.
func Fields(vals ...fieldValidation) fieldValidations {
	return fieldValidations(vals)
}

// fieldValidation is a struct that contains a set of rules
// that form values must comply with for a specific field.
type fieldValidation struct {
	Field      string
	Validators []ValidatorFn
}

type fieldValidations []fieldValidation

// Validate is the main method we will use to perform the validations on a form.
func (v fieldValidations) Validate(form url.Values) Errors {
	verrs := make(map[string][]error)

	for _, validation := range v {
		for _, rule := range validation.Validators {
			err := rule(form[validation.Field])
			if err == nil {
				continue
			}

			verrs[validation.Field] = append(verrs[validation.Field], err)
		}
	}

	return verrs
}

// Errors is a convenience field to map the form field name to the error message.
type Errors map[string][]error

func (ee Errors) HasAny() bool {
	return len(ee) > 0
}

func (ee Errors) Has(field string) bool {
	found, ok := ee[field]
	return ok && len(found) > 0
}

// ErrorStringFor returns the error string concatenating
// all the error messages for a specific field.
func (err Errors) ErrorStringFor(field string) string {
	if !err.Has(field) {
		return ""
	}
	parts := []string{}
	for _, v := range err[field] {
		parts = append(parts, v.Error())
	}

	return strings.Join(parts, ". ")
}

// ValidatorFn is a condition that must be satisfied by all values in a specific form field.
// Otherwise the rule will return an error
type ValidatorFn func([]string) error
