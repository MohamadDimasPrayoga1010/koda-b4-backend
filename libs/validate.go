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


func FormatProductValidationError(err error) map[string]string {
	errors := map[string]string{}

	for _, e := range err.(validator.ValidationErrors) {
		field := e.Field()
		tag := e.ActualTag()

		switch field {
		case "Title":
			errors["title"] = "Title wajib diisi"

		case "Description":
			if tag == "required" {
				errors["description"] = "Description wajib diisi"
			}
			if tag == "min" {
				errors["description"] = "Description minimal 10 karakter"
			}

		case "BasePrice":
			if tag == "required" {
				errors["basePrice"] = "Base Price wajib diisi"
			}
			if tag == "gt" {
				errors["basePrice"] = "Base Price harus lebih besar dari 0"
			}

		case "Stock":
			if tag == "required" {
				errors["stock"] = "Stock wajib diisi"
			}
			if tag == "gte" {
				errors["stock"] = "Stock tidak boleh kurang dari 0"
			}

		case "CategoryID":
			if tag == "required" {
				errors["categoryId"] = "Category ID wajib diisi"
			}
			if tag == "gt" {
				errors["categoryId"] = "Category ID harus lebih besar dari 0"
			}

		case "VariantID":
			if tag == "gt" {
				errors["variantId"] = "Variant ID harus angka valid dan > 0"
			}

		case "Sizes":
			if tag == "gt" {
				errors["sizes"] = "Size ID harus angka valid dan > 0"
			}

		case "Images":
			if tag == "required" {
				errors["images"] = "Minimal 1 gambar diperlukan"
			}
			if tag == "min" {
				errors["images"] = "Minimal 1 gambar diperlukan"
			}
		}
	}

	return errors
}
