---
title: Forms
index: 3
---
Leapkit ships with a form package that provides a flexible and reusable way to validate form data by defining a set of validation rules that can be applied to form fields.

## Validations
The `form/validate` package that offers a flexible and reusable way to validate form data by defining a set of validation rules that can be applied to form fields. Validations are a set of rules stablished for different fields passed in the request.

You can define these Validations to be used in your http handlers by and call the `form.Validate` function passing the `req` variable (*http.Request) and handling the `validate.Errors` returned. For example:

```go
// Validations for multiple fields
rules := validate.Fields(
	validate.Field("email", validate.Required("email is required")),
	validate.Field(
		"password",
		// Validations for the field
		validate.Required("password is required"),
		validate.MatchRegex(regexp.MustCompile(`^(?=.*[a-z])..`),"Password must be ..."),
	),
)

verrs := form.Validate(req, rules)
if verrs.HasAny() {
	 // handle validation errors...
}
```

When using leapkit validations consider the following:

### Fields
Unlike other libraries, in this package fields are taken from the request form. Names used for the validation must match the field that wants to be validated. It is expected that the form is already parsed when validating. The forms.Validate function does not parse the request form.

### Validations
A field can have multiple validations specified (Required, Length, Regex ...) and each validation can define an error message that will be returned if the validation does not pass. A field can use both built-in and custom validations.

### Errors
The output from the Validate function is a `validate.Errors` variable, which internally is a `map[string][]error` and provides some helpful functions. This structure allows to return multiple errors for a single field.

Errors has three useful methods:

- `HasAny()` returns true if there are any errors.
- `Has(field string)` returns true if there are errors for the given field.
- `StringFieldFor(field string)` returns a string with all the errors for the given field.

### Built-in Rules

You can build your set of rules for each validation by using the package's built-in functions.

```go
// General Rules:
func Required(message ...string) Rule

// String Rules:
func Matches(field string, message ...string) Rule
func MatchRegex(re *regexp.Regexp, message ...string) Rule
func MinLength(min int, message ...string) Rule
func MaxLength(max int, message ...string) Rule
func WithinOptions(options []string, message ...string) Rule

// Number Rules:
func EqualTo(value float64, message ...string) Rule
func LessThan(value float64, message ...string) Rule
func LessThanOrEqualTo(value float64, message ...string) Rule
func GreaterThan(value float64, message ...string) Rule
func GreaterThanOrEqualTo(value float64, message ...string) Rule

// Time Rules:
func TimeEqualTo(u time.Time, message ...string) Rule
func TimeBefore(u time.Time, message ...string) Rule
func TimeBeforeOrEqualTo(u time.Time, message ...string) Rule
func TimeAfter(u time.Time, message ...string) Rule
func TimeAfterOrEqualTo(u time.Time, message ...string) Rule

// Utility rules
func EmailValid(message ...string) Rule
func URLValid(message ...string) Rule
```

### Custom validation Rules

Alternatively, you can create your own validation functions. As long as these follow the `validate.ValidatorFn` (`func([]string) error`) signature you can apply these to fields. Like in the following example:

```go
// IsUnique checks in the database for a user email to be unique.
func UniqueEmail(db *sqlx.DB ) func([]string) error {
	return func(emails []string) error {
    query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
		stmt, err := db.Prepare(query)
		if err != nil {
			return err
		}

		for _, email := range emails {
			var exists bool
			if err := stmt.QueryRow(email).Scan(&exists); err != nil {
				return err
			}

			if exists {
				return fmt.Errorf("email '%s' already exists.", email)
			}
		}

		return nil
	}
}

// ...
rules := validation.Fields(
	// email field is checked to be:
	// - Present
	// - A valid email
	// - Unique in the database
	validation.Field("email", validate.EmailValid(), validate.Email(), UniqueEmail(db))
  ...
)
...
```
