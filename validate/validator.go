package validate

import (
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func LoadValidator(){
	Validate = validator.New(validator.WithRequiredStructEnabled())
}
func ErrorParse(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)
	for _, f := range verr {
		errs[f.Field()] = f.ActualTag()
	}
	return errs
}