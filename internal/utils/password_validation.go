package utils

import (
	"fmt"
	"regexp"
	"strings"
)

type PasswordValidator struct {
	minLen int
}

func NewPasswordValidator(minLen int) *PasswordValidator {
	if minLen < 8 {
		minLen = 8
	}
	return &PasswordValidator{minLen: minLen}
}

func (pv *PasswordValidator) Validate(password string) error {
	password = strings.TrimSpace(password)
	if len(password) < pv.minLen {
		return fmt.Errorf("password must be at least %d characters long", pv.minLen)
	}

	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Character class checks
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	var errors []string
	if !hasUpper {
		errors = append(errors, "at least one uppercase letter")
	}
	if !hasLower {
		errors = append(errors, "at least one lowercase letter")
	}
	if !hasDigit {
		errors = append(errors, "at least one digit")
	}
	if !hasSpecial {
		errors = append(errors, "at least one special character")
	}

	if len(errors) > 0 {
		return fmt.Errorf("password needs: %s", strings.Join(errors, ", "))
	}

	return nil
}
