package common

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	var ok bool
	validate, ok = binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}
	validate.RegisterValidation("test", func(fl validator.FieldLevel) bool {
		return true
	})
}

func Validator(data any) error {
	return validate.Struct(data)
}
