package rule

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func Regex(fl validator.FieldLevel) bool {
	// "^[a-zA-Z0-9]$"
	return regexp.MustCompile(fl.Param()).MatchString(fl.Field().String())
}

func Password(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
