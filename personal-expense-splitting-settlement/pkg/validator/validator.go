package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom password validation
	// Password must contain: min 8 chars, uppercase, lowercase, number, special char
	validate.RegisterValidation("password", validatePassword)
}

// validatePassword validates password strength
// Requirements: Min 8 chars, uppercase, lowercase, number, special char
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}

	// Check for uppercase
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// Check for lowercase
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	// Check for number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	// Check for special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// Validate validated a struct
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
