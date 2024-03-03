package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	v.FieldErrors[field] = message
}

func (v *Validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func WithinMaxChars(value string, maxChars int) bool {
	if utf8.RuneCountInString(strings.TrimSpace(value)) > maxChars {
		return false
	}
	return true
}

func Match[T comparable](firstValue, secondValue T) bool {
	if firstValue == secondValue {
		return true
	} else {
		return false
	}
}

func PermittedValue[T comparable](value T, permitted ...T) bool {
	for i := range permitted {
		if value == permitted[i] {
			return true
		}
	}
	return false
}

func MinChars(value string, minLength int) bool {
	if utf8.RuneCountInString(value) >= minLength {
		return true
	}
	return false
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}
