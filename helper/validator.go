package helper

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct based on validation tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// FormatValidationErrors formats validation errors into a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "email":
				errors[field] = field + " must be a valid email"
			case "min":
				errors[field] = field + " must be at least " + e.Param() + " characters"
			case "max":
				errors[field] = field + " must be at most " + e.Param() + " characters"
			case "oneof":
				errors[field] = field + " must be one of: " + e.Param()
			case "uuid":
				errors[field] = field + " must be a valid UUID"
			case "url":
				errors[field] = field + " must be a valid URL"
			default:
				errors[field] = field + " is invalid"
			}
		}
	}

	return errors
}
