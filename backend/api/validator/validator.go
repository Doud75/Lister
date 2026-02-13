package validator

import (
	"errors"
	"regexp"
	"unicode"

	"github.com/microcosm-cc/bluemonday"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)
	sanitizer     = bluemonday.StrictPolicy()
)

func ValidateUsername(username string) error {
	if !usernameRegex.MatchString(username) {
		return errors.New("le nom d'utilisateur doit contenir entre 3 et 50 caractères alphanumériques ou underscore")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("le mot de passe doit contenir au moins 8 caractères")
	}

	var hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	if !hasDigit {
		return errors.New("le mot de passe doit contenir au moins un chiffre")
	}
	if !hasSpecial {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}

	return nil
}

func Sanitize(input string) string {
	return sanitizer.Sanitize(input)
}
