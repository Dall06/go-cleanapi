// Package utils_test is a test package for utils
package utils_test

import (
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewString(test *testing.T) {
	successfulCases := []struct {
		name string
	}{
		{
			name: "it should generate an string",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := utils.NewUUIDGenerator()
			s := u.NewString()

			assert.NotEmpty(t, s, "expected string not to be empty")
		})
	}
}

func TestNewUUID(test *testing.T) {
	successfulCases := []struct {
		name string
	}{
		{
			name: "it should generate an uuid",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := utils.NewUUIDGenerator()
			s := u.NewUUID()

			assert.IsType(t, uuid.UUID{}, s, "expecting an uuid but got other")
			assert.NotEmpty(t, s, "expected an uuid not to be empty")
		})
	}
}
