package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"log"
	"unicode"
)

var (
	validations = []ValidationRule{
		{Validate: "required", Translator: "{0} is a required field"},
		{Validate: "email", Translator: "{0} must be a valid email address"},
		{Validate: "passwd", Translator: "{0} must be 8 or more characters and contain at least one of each uppercase, lowercase, number and symbol"},
		{Validate: "uuid", Translator: "{0} must be a valid uuid"},
		{Validate: "requiredUuid", Translator: "{0} must be a valid uuid and must not be nil"},
	}
)

type ValidationRule struct {
	Validate   string
	Translator string
}

type Validator struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

// New creates and returns a new validator
//TODO: needs to be more dynamic as the hardcoded validations at the bottom are not ideal.
func New() *Validator {
	log.Println("Setting up translation and validation")
	v := &Validator{}

	enTranslator := en.New()

	uni := ut.New(enTranslator, enTranslator)

	var found bool
	v.Translator, found = uni.GetTranslator("en")

	if !found {
		log.Println("translator not found")
		return nil
	}

	v.Validate = validator.New()

	if err := en2.RegisterDefaultTranslations(v.Validate, v.Translator); err != nil {
		log.Printf("error registering default translations: %s", err.Error())
		return nil
	}
	for _, val := range validations {
		if v.addTranslator(val.Validate, val.Translator, v.Translator) != nil {
			return nil
		}
	}

	if v.addValidation("passwd", passwordValidator) != nil {
		return nil
	}
	if v.addValidation("uuid", uuidValidator) != nil {
		return nil
	}
	if v.addValidation("requiredUuid", requiredUUIDValidator) != nil {
		return nil
	}
	return v
}

func (v *Validator) addTranslator(name, message string, translator ut.Translator) error {
	err := v.Validate.RegisterTranslation(name, translator, func(ut ut.Translator) error {
		return ut.Add(name, message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(name, fe.Field())
		return t
	})

	if err != nil {
		log.Printf("error registering default translations: %s", err.Error())
		return err
	}
	return nil
}

func (v *Validator) addValidation(name string, validationAction func(fl validator.FieldLevel) bool) error {
	return v.Validate.RegisterValidation(name, validationAction)
}

func passwordValidator(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) < 8 {
		return false
	}
	inputPassword := fl.Field().String()

	totalLength := 0
	containsLowerCase := false
	containsUppercase := false
	containsNumbers := false
	containsSymbols := false

	for _, c := range inputPassword {
		switch {
		case unicode.IsNumber(c):
			containsNumbers = true

		case unicode.IsLower(c):
			containsLowerCase = true

		case unicode.IsUpper(c):
			containsUppercase = true

		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			containsSymbols = true
		}

		totalLength++
	}

	if totalLength < 8 || !containsNumbers || !containsSymbols || !containsUppercase || !containsLowerCase {
		return false
	}

	return true
}

func uuidValidator(fl validator.FieldLevel) bool {
	_, err := uuid.Parse(fl.Field().String())
	if err != nil {
		return false
	}
	return true
}

func requiredUUIDValidator(fl validator.FieldLevel) bool {
	uid, err := uuid.Parse(fl.Field().String())
	if err != nil {
		return false
	}
	if uid == uuid.Nil {
		return false
	}
	return true
}
