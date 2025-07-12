package validator

import (
	"fmt"
	"sort"
	"strings"
)

type ValidationError struct {
	fields []FieldError
}

func NewValidationError(fields []FieldError) *ValidationError {
	return &ValidationError{fields: fields}
}

func (e *ValidationError) Fields() []FieldError {
	return e.fields
}

func (e *ValidationError) Error() string {
	if len(e.fields) == 0 {
		return "validation failed"
	}

	var reasons []string
	for _, f := range e.fields {
		reasons = append(reasons, f.Reason())
	}
	return fmt.Sprintf("validation failed: %s", strings.Join(reasons, ", "))
}

func (e *ValidationError) HasErrors() bool {
	return len(e.fields) > 0
}

func (e *ValidationError) SortByField() *ValidationError {
	sorted := make([]FieldError, len(e.fields))
	copy(sorted, e.fields)

	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Field() != sorted[j].Field() {
			return sorted[i].Field() < sorted[j].Field()
		}
		return sorted[i].Code() < sorted[j].Code()
	})

	return &ValidationError{fields: sorted}
}

type FieldError struct {
	field  string
	code   string
	reason string
}

func NewFieldError(field, code, reason string) FieldError {
	return FieldError{
		field:  field,
		code:   code,
		reason: reason,
	}
}

func (e FieldError) Field() string { return e.field }

func (e FieldError) Code() string { return e.code }

func (e FieldError) Reason() string { return e.reason }

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.field, e.reason)
}

type ValidationResult struct {
	errors *ValidationError
}

func Success() ValidationResult {
	return ValidationResult{errors: nil}
}

func Failure(fields []FieldError) ValidationResult {
	return ValidationResult{errors: NewValidationError(fields)}
}

func (r ValidationResult) IsValid() bool {
	return r.errors == nil || !r.errors.HasErrors()
}

func (r ValidationResult) Errors() *ValidationError {
	return r.errors
}

func (r ValidationResult) MustBeValid() {
	if !r.IsValid() {
		panic(fmt.Sprintf("validation must pass but got: %v", r.errors))
	}
}
