package validate_test

import (
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/leapkit/leapkit/core/form/validate"
)

func TestRuleRequired(test *testing.T) {
	// Given a form with not-empty field values, Then the validate.Required rule should return no error.
	test.Run("correct form has field values", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.Required()),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form without field, Then the validate.Required rule should return error.
	test.Run("incorrect form does not have field", func(t *testing.T) {
		form := url.Values{}

		validator := validate.Fields(
			validate.Field("input_field", validate.Required()),
		)

		verrs := validator.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with at least one empty field value, Then the validate.Required rule should return error
	test.Run("incorrect form field has at least one empty value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1", "", "value_3"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.Required()),
		)

		verrs := validations.Validate(form)

		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleMatches(test *testing.T) {
	// Given a form with values that match the field, Then the Matches rule should return no error.
	test.Run("correct form field values match with field", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.Matches("value_1")),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values that don't match the field, Then the Matches rule should return error.
	test.Run("incorrect form field values do not match with field", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.Matches("value_2")),
		)

		verrs := validations.Validate(form)

		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleMatchRegex(test *testing.T) {
	// Given a form with values that match with the regular expression, Then the MatchRegex rule should return no error.
	test.Run("correct form field values match with the regular expression", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"seafood"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MatchRegex(regexp.MustCompile(`foo.*`))),
		)

		verrs := validations.Validate(form)

		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values that don't match with the regular expression, Then the MatchRegex rule should return error.
	test.Run("incorrect form field values do not match with the regular expression", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"seafood"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MatchRegex(regexp.MustCompile(`bar.*`))),
		)

		verrs := validations.Validate(form)

		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleEqualTo(test *testing.T) {
	// Given a form with values less than compared value, Then the EqualTo rule should return no error.
	test.Run("correct form field value is equal to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10.36"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.EqualTo(10.36)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values equal to compared value, Then the EqualTo rule should return error.
	test.Run("incorrect form field value is different to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.EqualTo(20)),
		)

		verrs := validations.Validate(form)

		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with no number values, Then the EqualTo rule should return error.
	test.Run("incorrect form field value is not a number", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"invalid value"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.EqualTo(5), validate.Required()),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleLessThan(test *testing.T) {
	// Given a form with values less than compared value, Then the LessThan rule should return no error.
	test.Run("correct form field value is less to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThan(20)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values equal to compared value, Then the LessThan rule should return error.
	test.Run("incorrect form field value is equal to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThan(10)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with values greater than compared value, Then the LessThan rule should return error.
	test.Run("incorrect form field value is greater than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThan(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with no number values, Then the LessThan rule should return error.
	test.Run("incorrect form field value is not a number", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"invalid value"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThan(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleLessThanOrEqualTo(test *testing.T) {
	// Given a form with values less than compared value, Then the LessThanOrEqualTo rule should return no error.
	test.Run("correct form field value is less to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThanOrEqualTo(20)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values equal to compared value, Then the LessThanOrEqualTo rule should return no error.
	test.Run("correct form field value is equal to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThanOrEqualTo(10)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values greater than compared value, Then the LessThanOrEqualTo rule should return error.
	test.Run("incorrect form field value is greater than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThanOrEqualTo(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with no number values, Then the LessThanOrEqualTo rule should return error.
	test.Run("incorrect form field value is not a number", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"invalid value"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.LessThanOrEqualTo(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleGreaterThan(test *testing.T) {
	// Given a form with values greater than compared value, Then the GreaterThan rule should return no error.
	test.Run("correct form field value is greater than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThan(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values equal to compared value, Then the GreaterThan rule should return error.
	test.Run("incorrect form field value is equal to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThan(10)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with values less than compared value, Then the GreaterThan rule should return error.
	test.Run("incorrect form field value is less than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThan(20)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with no number values, Then the GreaterThan rule should return error.
	test.Run("incorrect form field value is not a number", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"invalid value"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThan(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleGreaterThanOrEqualTo(test *testing.T) {
	// Given a form with values greater than compared value, Then the GreaterThanOrEqualTo rule should return no error.
	test.Run("correct form field value is greater than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThanOrEqualTo(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values equal to compared value, Then the GreaterThanOrEqualTo rule should return no error.
	test.Run("correct form field value is equal to compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThanOrEqualTo(10)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form with values less than compared value, Then the GreaterThanOrEqualTo rule should return error.
	test.Run("incorrect form field value is less than compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"10"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThanOrEqualTo(20)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form with no number values, Then the GreaterThanOrEqualTo rule should return error.
	test.Run("incorrect form field value is not a number", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"invalid value"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.GreaterThanOrEqualTo(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleMinLength(test *testing.T) {
	// Given a form field values with a length greater than the compared value, Then the MinLength rule should return no error.
	test.Run("correct form field values with a length greater than the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lorem ipsum"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MinLength(3)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values with a length equal to the compared value, Then the MinLength rule should return no error.
	test.Run("correct form field values with a length equal to the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lorem ipsum"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MinLength(11)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values with a length less than the compared value, Then the MinLength rule should return error.
	test.Run("incorrect form field values with a length less than the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lo"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MinLength(11)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleMaxLength(test *testing.T) {
	// Given a form field values with a length less than the compared value, Then the MaxLength rule should return no error.
	test.Run("correct form field values with a length greater than the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lorem ipsum"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MaxLength(20)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values with a length equal to the compared value, Then the MaxLength rule should return no error.
	test.Run("correct form field values with a length equal to the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lorem ipsum"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MaxLength(11)),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values with a length greater than the compared value, Then the MaxLength rule should return error.
	test.Run("incorrect form field values with a length less than the compared value", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"lorem ipsum"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.MaxLength(5)),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleWithinOptions(test *testing.T) {
	// Given a form field with values that are in the option list, Then the WithinOptions rule should return no error.
	test.Run("correct form field values are in the option list", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1", "value_2"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.WithinOptions([]string{"value_1", "value_2", "value_3"})),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field with at leas a value that is not in the option list, Then the WithinOptions rule should return error.
	test.Run("incorrect a form field value is not in the option list", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"value_1", "value_4"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.WithinOptions([]string{"value_1", "value_2", "value_3"})),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleTimeEqualTo(test *testing.T) {
	// Given a form field values that are times equal to the compared time, Then the TimeEqualTo rule should return no error.
	test.Run("correct form field values are times equal to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeEqualTo(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times different to the compared time, Then the TimeEqualTo rule should return error.
	test.Run("incorrect form field values are times different to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are not times, Then the TimeEqualTo rule should return error.
	test.Run("incorrect form field values that are not times", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"is not a time"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeEqualTo(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleTimeBefore(test *testing.T) {
	// Given a form field values that are times before to the compared time, Then the TimeBefore rule should return no error.
	test.Run("correct form field values are times before to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBefore(time.Date(2028, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times equal to the compared time, Then the TimeBefore rule should return error.
	test.Run("incorrect form field values are times equal to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBefore(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are times after to the compared time, Then the TimeBefore rule should return error.
	test.Run("incorrect form field values are times after to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBeforeOrEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are not times, Then the TimeBefore rule should return error.
	test.Run("incorrect form field values are not times", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"is not a time"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBefore(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleTimeBeforeOrEqualTo(test *testing.T) {
	// Given a form field values that are times before to the compared time, Then the TimeBeforeOrEqualTo rule should return no error.
	test.Run("correct form field values are times before to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBeforeOrEqualTo(time.Date(2028, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times equal to the compared time, Then the TimeBeforeOrEqualTo rule should return no error.
	test.Run("correct form field values are times equal to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBeforeOrEqualTo(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times after to the compared time, Then the TimeBeforeOrEqualTo rule should return error.
	test.Run("incorrect form field values are times after to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBeforeOrEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are not times, Then the TimeBeforeOrEqualTo rule should return error.
	test.Run("incorrect form field values are not times", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"is not a time"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeBeforeOrEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleTimeAfter(t *testing.T) {
	// Given a form field values that are times after to the compared time, Then the TimeAfter rule should return no error.
	t.Run("correct form field values are times after to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfter(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times equal to the compared time, Then the TimeAfter rule should return error.
	t.Run("incorrect form field values are times equal to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfter(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are times before to the compared time, Then the TimeAfter rule should return error.
	t.Run("incorrect form field values are times before to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfter(time.Date(2028, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are not times, Then the TimeAfter rule should return error.
	t.Run("incorrect form field values are not times", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"is not a time"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfter(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestRuleTimeAfterOrEqualTo(test *testing.T) {
	// Given a form field values that are times after to the compared time, Then the TimeAfterOrEqualTo rule should return no error.
	test.Run("correct form field values are times after to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfterOrEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times equal to the compared time, Then the TimeAfterOrEqualTo rule should return no error.
	test.Run("correct form field values are times equal to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfterOrEqualTo(time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	// Given a form field values that are times before to the compared time, Then the TimeAfterOrEqualTo rule should return error.
	test.Run("incorrect form field values are times before to the compared time", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"2026-06-26"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfterOrEqualTo(time.Date(2028, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})

	// Given a form field values that are not times, Then the TimeAfterOrEqualTo rule should return error.
	test.Run("incorrect form field values are not times", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"is not a time"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.TimeAfterOrEqualTo(time.Date(2025, time.June, 26, 0, 0, 0, 0, time.UTC))),
		)

		verrs := validations.Validate(form)
		if len(verrs) == 0 {
			t.Fatalf("verrs should have errors. verrs=%v", verrs)
		}
	})
}

func TestEmailValidRule(t *testing.T) {
	t.Run("correct addresses", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"a@pagano.id"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.EmailValid()),
		)

		verrs := validations.Validate(form)
		if len(verrs) != 0 {
			t.Fatalf("verrs should not have errors. verrs=%v", verrs)
		}
	})

	t.Run("incorrect addresses", func(t *testing.T) {
		tcases := map[string]string{
			"missing tld":                       "a@pagano",
			"missing domain":                    "a@.id",
			"missing username":                  "@pagano.id",
			"missing username and domain":       "@.id",
			"missing username and tld":          "@pagano.",
			"missing domain and tld":            "a@.",
			"missing username, domain, and tld": "@",
		}

		for name, email := range tcases {
			t.Run(name, func(t *testing.T) {
				form := url.Values{
					"input_field": []string{email},
				}

				validations := validate.Fields(
					validate.Field("input_field", validate.EmailValid()),
				)

				verrs := validations.Validate(form)
				if len(verrs) == 0 {
					t.Fatalf("verrs should have errors. verrs=%v", verrs)
				}
			})
		}
	})
}

func TestRuleURLValid(t *testing.T) {
	t.Run("correct URLs", func(t *testing.T) {
		form := url.Values{
			"input_field": []string{"https://wawand.co"},
		}

		validations := validate.Fields(
			validate.Field("input_field", validate.URLValid()),
		)

		verrs := validations.Validate(form)
		if len(verrs) > 0 {
			t.Fatalf("verrs must not have errors, verrs=%v", verrs)
		}
	})

	t.Run("incorrect URLs", func(t *testing.T) {
		tcases := map[string]string{
			"missing scheme": "wawand.co",
			"missing domain": "hssssss",
		}

		for name, u := range tcases {
			t.Run(name, func(t *testing.T) {
				form := url.Values{
					"input_field": []string{u},
				}

				validations := validate.Fields(
					validate.Field("input_field", validate.URLValid()),
				)

				verrs := validations.Validate(form)
				if len(verrs) == 0 {
					t.Fatalf("verrs should have errors. verrs=%v", verrs)
				}
			})
		}
	})
}
