package utils

import "github.com/go-playground/validator/v10"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationError(err error) []FieldError {
	var errors []FieldError

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, FieldError{
			Field:   e.Field(),
			Message: e.Tag(), // bisa di-custom nanti
		})
	}

	return errors
}
