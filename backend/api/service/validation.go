package service

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return errors.New("le nom d'utilisateur doit contenir au moins 3 caractères")
	}
	if len(username) > 50 {
		return errors.New("le nom d'utilisateur ne peut pas dépasser 50 caractères")
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("le nom d'utilisateur ne peut contenir que des lettres, des chiffres et des underscores")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("le mot de passe doit contenir au moins 8 caractères")
	}
	var (
		hasUpper   bool
		hasNumber  bool
		hasSpecial bool
	)
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	if !hasUpper {
		return errors.New("le mot de passe doit contenir au moins une majuscule")
	}
	if !hasNumber {
		return errors.New("le mot de passe doit contenir au moins un chiffre")
	}
	if !hasSpecial {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}
	return nil
}

func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}
