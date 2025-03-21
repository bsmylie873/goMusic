package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("validSex", validateSex)
	validate.RegisterValidation("validTitle", validateTitle)
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func validateSex(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	return value >= 1 && value <= 3
}

func validateTitle(fl validator.FieldLevel) bool {
	value := fl.Field().Int()
	return value >= 1 && value <= 4
}
