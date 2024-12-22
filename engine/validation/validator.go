package validation

import (
"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validator *validator.Validate
}

func NewValidation() Validator {
	return Validator{
		Validator: validator.New(),
	}
}

func (ths *Validator) Validate(data interface{}) map[string]interface{} {
	if err := ths.Validator.Struct(data); err != nil {
		errorMap := make(map[string]interface{})
		switch v := err.(type) {
		case *validator.InvalidValidationError:
			errorMap[v.Type.Name()] = v.Error()
		case validator.ValidationErrors:
			for _, value := range v {
				errorMap[value.Namespace()] = value.Tag() + " " + value.Param()
			}
		default:
			errorMap["message"] = err.Error()
		}

		return errorMap
	}
	return nil
}
