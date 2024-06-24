package validate

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

// Required function validates the form field has no-empty values.
func Required(message ...string) ValidatorFn {
	return func(values []string) error {
		hasEmptyValues := slices.ContainsFunc(values, func(val string) bool {
			return strings.TrimSpace(val) == ""
		})

		if len(values) > 0 && !hasEmptyValues {
			return nil
		}

		return newError("This field is required.", message...)
	}
}

// Match function validates the form field values with a string.
func Matches(field string, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if val == field {
				continue
			}

			return newError(fmt.Sprintf("'%s' does not match with '%s'.", val, field), message...)
		}

		return nil
	}
}

// MatchRegex function validates the form field values with a regular expression.
func MatchRegex(re *regexp.Regexp, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if re.MatchString(val) {
				continue
			}

			return newError(fmt.Sprintf("'%s' does not match with '%s'.", val, re), message...)
		}

		return nil
	}
}

// EqualTo function validates that field values are equal to a compared value.
func EqualTo(value float64, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("is not a number")
			}

			if n == value {
				continue
			}

			return newError(fmt.Sprintf("%s must be equal to than %f.", val, value), message...)
		}

		return nil
	}
}

// LessThan function validates that the field values are less than a value.
func LessThan(value float64, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("is not a number")
			}

			if n < value {
				continue
			}

			return newError(fmt.Sprintf("%s must be less than %f.", val, value), message...)
		}

		return nil
	}
}

// LessThanOrEqualTo function validates that the field values are less than or equal to a value.
func LessThanOrEqualTo(value float64, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("is not a number")
			}

			if n <= value {
				continue
			}

			return newError(fmt.Sprintf("%s must be less than or equal to %f.", val, value), message...)
		}

		return nil
	}
}

// GreaterThan function validates that the field values are greater than a value.
func GreaterThan(value float64, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("is not a number")
			}

			if n > value {
				continue
			}

			return newError(fmt.Sprintf("%s must be greater than %f.", val, value), message...)
		}

		return nil
	}
}

// GreaterThanOrEqualTo function validates that the field values are greater than or equal to a value.
func GreaterThanOrEqualTo(value float64, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			n, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return errors.New("is not a number")
			}

			if n >= value {
				continue
			}

			return newError(fmt.Sprintf("%s must be greater than or equal to %f.", val, value), message...)
		}

		return nil
	}
}

// MinLength function validates that the values' lengths are greater than or equal to min.
func MinLength(min int, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if len(strings.TrimSpace(val)) >= min {
				continue
			}

			return newError(fmt.Sprintf("'%s' must not be less than %d characters.", val, min), message...)
		}

		return nil
	}
}

// MaxLength function validates that the values' lengths are less than or equal to max.
func MaxLength(max int, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if len(strings.TrimSpace(val)) <= max {
				continue
			}

			return newError(fmt.Sprintf("'%s' must not exceed %d characters.", val, max), message...)
		}

		return nil
	}
}

// WithinOptions function validates that values are in the option list.
func WithinOptions(options []string, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if slices.Contains(options, val) {
				continue
			}

			return newError(fmt.Sprintf("'%s' is not in the options.", val), message...)

		}

		return nil
	}
}

// ValidUUID function validates that the values are valid UUIDs.
func ValidUUID(message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			if !uuid.FromStringOrNil(val).IsNil() {
				continue
			}

			return newError(fmt.Sprintf("'%s' is not a valid uuid.", val), message...)
		}

		return nil
	}
}

// TimeEqualTo function validates that the values are equal an specific time.
func TimeEqualTo(u time.Time, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, value := range values {
			t, err := parseTime(value)
			if err != nil {
				return errors.New("is not a time")
			}

			if t.Equal(u) {
				continue
			}

			return newError(fmt.Sprintf("Time should be equal to '%s'.", u.Format(time.DateOnly)), message...)
		}

		return nil
	}
}

// TimeBefore function validates that the values are before an specific time.
func TimeBefore(u time.Time, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, value := range values {
			t, err := parseTime(value)
			if err != nil {
				return errors.New("is not a time")
			}

			if t.Before(u) {
				continue
			}

			return newError(fmt.Sprintf("Time should be before than '%s'.", u.Format(time.DateOnly)), message...)
		}

		return nil
	}
}

// TimeBeforeOrEqualTo function validates that the values are before or equal to an specific time.
func TimeBeforeOrEqualTo(u time.Time, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, value := range values {
			t, err := parseTime(value)
			if err != nil {
				return errors.New("is not a time")
			}

			if t.Before(u) || t.Equal(u) {
				continue
			}

			return newError(fmt.Sprintf("Time should be before or equal to '%s'.", u.Format(time.DateOnly)), message...)
		}

		return nil
	}
}

// TimeAfter function validates that the values are after an specific time.
func TimeAfter(u time.Time, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			t, err := parseTime(val)
			if err != nil {
				return newError("invalid time", message...)
			}

			if t.After(u) {
				continue
			}

			return newError(fmt.Sprintf("Time should be after '%s'.", u.Format(time.DateOnly)), message...)
		}

		return nil
	}
}

// TimeAfterOrEqualTo function validates that the values are after or equal to an specific time.
func TimeAfterOrEqualTo(u time.Time, message ...string) ValidatorFn {
	return func(values []string) error {
		for _, val := range values {
			t, err := parseTime(val)
			if err != nil {
				return newError("invalid time", message...)
			}

			if t.After(u) || t.Equal(u) {
				continue
			}

			return newError(fmt.Sprintf("Time should be after or equal to '%s'.", u.Format(time.DateOnly)), message...)
		}

		return nil
	}
}

func parseTime(strTime string) (time.Time, error) {
	layouts := []string{
		time.DateOnly,
		time.Layout,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
		time.DateTime,
		time.TimeOnly,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, strTime)
		if err != nil {
			continue
		}

		return t, nil
	}

	return time.Time{}, errors.New("invalid time")
}

func newError(str string, override ...string) error {
	override = append(override, str)

	return errors.New(cmp.Or(override...))
}
