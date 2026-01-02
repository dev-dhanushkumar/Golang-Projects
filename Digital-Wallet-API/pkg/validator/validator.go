package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct
func Validate(data interface{}) error {
	return validate.Struct(data)
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}
