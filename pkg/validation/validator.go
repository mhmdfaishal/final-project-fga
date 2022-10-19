package validation

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func Validate(i interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(i)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) { // loop through all the errors
			var element ErrorResponse                   // create a new error element
			element.FailedField = err.StructNamespace() // get the field name
			element.Tag = err.Tag()                     // get the tag
			element.Value = err.Param()                 // get the value
			errors = append(errors, &element)           // append the error element to the errors array
		}
	}

	return errors
}
