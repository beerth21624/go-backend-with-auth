package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Rule[T any] func(T) error

func Validate[T any](v T, rules ...Rule[T]) error {
	for _, r := range rules {
		if err := r(v); err != nil {
			return err
		}
	}
	return nil
}

func Collect[T any](v T, rules ...Rule[T]) []error {
	var errs []error
	for _, r := range rules {
		if err := r(v); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func All[T any](rules ...Rule[T]) Rule[T] {
	return func(v T) error { return Validate(v, rules...) }
}

func Any[T any](rules ...Rule[T]) Rule[T] {
	return func(v T) error {
		var last error
		for _, r := range rules {
			if err := r(v); err == nil {
				return nil
			} else {
				last = err
			}
		}
		return last
	}
}

func Not[T any](rule Rule[T], msg string) Rule[T] {
	return func(v T) error {
		if err := rule(v); err == nil {
			return errors.New(msg)
		}
		return nil
	}
}

func NotEmpty(msg string) Rule[string] {
	if msg == "" {
		msg = "must not be empty"
	}
	return func(s string) error {
		if s == "" {
			return errors.New(msg)
		}
		return nil
	}
}

func MaxLen(msg string, max int) Rule[string] {
	return func(s string) error {
		if len(s) > max {
			return errors.New(msg)
		}
		return nil
	}
}

func MinLen(msg string, min int) Rule[string] {
	return func(s string) error {
		if len(s) < min {
			return errors.New(msg)
		}
		return nil
	}
}

func Min[T ~int | ~int64 | ~float64](msg string, minVal T) Rule[T] {
	return func(v T) error {
		if v < minVal {
			return errors.New(msg)
		}
		return nil
	}
}

func Max[T ~int | ~int64 | ~float64](msg string, maxVal T) Rule[T] {
	return func(v T) error {
		if v > maxVal {
			return errors.New(msg)
		}
		return nil
	}
}

func Contains(msg string, substr string) Rule[string] {
	return func(s string) error {
		if !strings.Contains(s, substr) {
			return errors.New(msg)
		}
		return nil
	}
}

func Match(msg string, pattern string) Rule[string] {
	return func(s string) error {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			return errors.New("invalid pattern")
		}
		if !matched {
			return errors.New(msg)
		}
		return nil
	}
}

func OneOf[T comparable](msg string, values ...T) Rule[T] {
	return func(v T) error {
		for _, val := range values {
			if v == val {
				return nil
			}
		}
		return errors.New(msg)
	}
}

func Range[T ~int | ~int64 | ~float64](msg string, min, max T) Rule[T] {
	return func(v T) error {
		if v < min || v > max {
			return errors.New(msg)
		}
		return nil
	}
}

func Field[S any, F any](extract func(S) F, rule Rule[F]) Rule[S] {
	return func(s S) error {
		return rule(extract(s))
	}
}

func CollectFieldErrors[T any](v T, rules ...FieldRule[T]) ValidationResult {
	var fieldErrors []FieldError
	for _, r := range rules {
		if err := r(v); err != nil {
			if fe, ok := err.(FieldError); ok {
				fieldErrors = append(fieldErrors, fe)
			} else {
				// fallback: แปลง error ธรรมดาเป็น FieldError
				fieldErrors = append(fieldErrors, NewFieldError("unknown", "error", err.Error()))
			}
		}
	}

	if len(fieldErrors) == 0 {
		return Success()
	}
	return Failure(fieldErrors)
}

type FieldRule[T any] func(T) error

func FieldValidation[S any, F any](fieldName string, extract func(S) F, rules ...Rule[F]) FieldRule[S] {
	return func(s S) error {
		fieldValue := extract(s)
		for _, rule := range rules {
			if err := rule(fieldValue); err != nil {
				// แปลง error เป็น FieldError พร้อม field name
				return NewFieldError(fieldName, "validation_failed", err.Error())
			}
		}
		return nil
	}
}

func ConditionalValidation[S any](condition func(S) bool, rule FieldRule[S]) FieldRule[S] {
	return func(s S) error {
		if condition(s) {
			return rule(s)
		}
		return nil
	}
}

func RequiredField(fieldName string, extract func(s interface{}) string) FieldRule[interface{}] {
	return func(s interface{}) error {
		value := extract(s)
		if strings.TrimSpace(value) == "" {
			return NewFieldError(fieldName, "required", fieldName+" is required")
		}
		return nil
	}
}

func MinLengthField(fieldName string, minLen int, extract func(s interface{}) string) FieldRule[interface{}] {
	return func(s interface{}) error {
		value := extract(s)
		if len(value) < minLen {
			return NewFieldError(fieldName, "min_length", fmt.Sprintf("%s must be at least %d characters", fieldName, minLen))
		}
		return nil
	}
}

func MaxLengthField(fieldName string, maxLen int, extract func(s interface{}) string) FieldRule[interface{}] {
	return func(s interface{}) error {
		value := extract(s)
		if len(value) > maxLen {
			return NewFieldError(fieldName, "max_length", fmt.Sprintf("%s must not exceed %d characters", fieldName, maxLen))
		}
		return nil
	}
}
