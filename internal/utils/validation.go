package utils

import "github.com/go-playground/validator/v10"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func FormatValidationError(err error) []FieldError {
	if err == nil {
		return nil
	}

	var errors []FieldError
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return []FieldError{
			{
				Field:   "request",
				Message: err.Error(),
			},
		}
	}

	for _, e := range validationErrors {
		errors = append(errors, FieldError{
			Field:   e.Field(),
			Message: e.Tag(), // bisa di-custom nanti
		})
	}

	return errors
}
