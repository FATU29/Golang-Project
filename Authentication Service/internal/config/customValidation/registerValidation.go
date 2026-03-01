package customValidation

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func register(v *validator.Validate, ok bool, key string, fnc validator.Func) {
	if ok {
		err := v.RegisterValidation(key, fnc)
		if err != nil {
			panic(err)
		}
	}

}

func RegisterCustomValidations() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	register(v, ok, "secure_password", PasswordValidation)
	register(v, ok, "email", EmailValidation)
}
