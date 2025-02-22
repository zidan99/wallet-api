package validator

import (
	"context"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v "github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

var Ctx context.Context

// Validate validates obj based on defined rulesets specified in the struct tag.
func Validate(obj any) []string {
	var errs []string

	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	validate := v.New()

	translations.RegisterDefaultTranslations(validate, trans)
	registerCustomTranslations(validate, trans)

	if err := validate.Struct(obj); err != nil {
		ve, ok := err.(v.ValidationErrors)
		if !ok {
			return make([]string, 0)
		}

		for _, e := range ve {
			message := e.Translate(trans)
			errs = append(errs, message)
		}
	}

	return errs
}

func registerCustomTranslations(validate *v.Validate, trans ut.Translator) {
	validate.RegisterTranslation(
		"hostname",
		trans,
		func(ut ut.Translator) error {
			return ut.Add("hostname", "{0} is not a valid hostname", true)
		},
		func(ut ut.Translator, fe v.FieldError) string {
			t, _ := ut.T("hostname", fe.Value().(string))
			return t
		},
	)
}
