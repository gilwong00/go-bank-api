package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUserName = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidateName  = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(v string, minLength int, maxLength int) error {
	l := len(v)
	if l < minLength || l > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUserName(value) {
		return fmt.Errorf("username must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidateName(value) {
		return fmt.Errorf("name must contain only letters or spaces")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}
