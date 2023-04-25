package utils_test

import (
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmail(test *testing.T) {
	successfulCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "it should check ok, is an email",
			input:    "test@test.com",
			expected: true,
		},
	}

	failedCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "it should not check ok, is not an email",
			input:    "testtest.com",
			expected: false,
		},
		{
			name:     "it should not check ok, is empty",
			input:    "",
			expected: false,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			v := utils.NewValidations()
			res := v.IsEmail(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			v := utils.NewValidations()
			res := v.IsEmail(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestIsPhone(test *testing.T) {
	successfulCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "it should check ok, is a phone",
			input:    "+991234567890",
			expected: true,
		},
	}

	failedCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "it should not check ok, is not a random string",
			input:    "dawedawedawe",
			expected: false,
		},
		{
			name:     "it should not check ok, is empty",
			input:    "",
			expected: false,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			v := utils.NewValidations()
			res := v.IsPhone(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			v := utils.NewValidations()
			res := v.IsPhone(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}
}
