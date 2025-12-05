package libs

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func FormatValidationError(err error) map[string]string {
	errors := map[string]string{}

	ve, ok := err.(validator.ValidationErrors)
	if !ok {
		errors["error"] = "Invalid input"
		return errors
	}

	for _, e := range ve {
		field := strings.ToLower(e.Field())
		tag := e.ActualTag()

		switch field {
		case "fullname":
			if tag == "required" {
				errors[field] = "Fullname wajib diisi"
			}

		case "email":
			if tag == "required" {
				errors[field] = "Email wajib diisi"
			}
			if tag == "email" {
				errors[field] = "Format email tidak valid"
			}

		case "password":
			if tag == "required" {
				errors[field] = "Password wajib diisi"
			}
			if tag == "min" {
				errors[field] = "Password minimal 6 karakter"
			}
		}
	}

	return errors
}
