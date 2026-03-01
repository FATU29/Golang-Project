package customValidation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func EmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// Simple regex for email validation
	var emailRegex = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
