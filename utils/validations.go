// Package utils is a package that provides general method for the api usage
package utils

import "regexp"

const (
	phoneRegexStr = `^\+[0-9]{10,}$`
	emailRegexStr = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

// Validations is an interface that extends validations
type Validations interface {
	IsPhone(string) bool
	IsEmail(string) bool
}

type validations struct{}

var _ Validations = (*validations)(nil)

// NewValidations is a constructor for validations
func NewValidations() Validations {
	return &validations{}
}

func (*validations) IsPhone(value string) bool {
	// use regex to check if it's a valid phone number
	phoneRegex := regexp.MustCompile(phoneRegexStr)
	r := phoneRegex.MatchString(value)
	return r
}

func (*validations) IsEmail(value string) bool {
	// use regex to check if it's a valid email address
	emailRegex := regexp.MustCompile(emailRegexStr)
	r := emailRegex.MatchString(value)
	return r
}
