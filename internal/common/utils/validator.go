package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()

	// register custom validation
	Validate.RegisterValidation("employeeid", EmployeeIDValidation)
}

func EmployeeIDValidation(fl validator.FieldLevel) bool {
	pattern := `^[0-9]+@[a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, fl.Field().String())
	return matched
}
