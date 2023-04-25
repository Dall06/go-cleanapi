// Package utils_test is a test package for utils
package utils_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUserJWT(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should create a jwt user based",
			input:  "im an id",
		},
	}

	failedCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should fail create a jwt user based, empty string",
			input:  "",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateUserJWT(tc.input)

			assert.NoError(t, err)                                              // it should be empty
			assert.NotEmpty(t, r, "expected a non-empty string, but got empty") // it should not be empty
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateUserJWT(tc.input)

			assert.Error(t, err)
			assert.Empty(t, r, "expected an empty string, but got something on it")
		})
	}
}

func TestCheckUserJWT(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should check a jwt user based",
			input:  "im an id",
		},
	}

	failedCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should fail check a jwt user based, empty string",
			input:  "",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateUserJWT(tc.input)

			if err != nil {
				t.Error("expected no error, but got:", err)
			}
			if r == "" {
				t.Error("expected a non-empty string, but got empty")
			}

			ok, err := jwt.CheckUserJwt(r)
			assert.NoError(t, err)
			assert.Equal(t, true, ok, "expected an ok, but got false")
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			ok, err := jwt.CheckUserJwt("")
			assert.Error(t, err)
			assert.Equal(t, false, ok, "expected false, but got true")
		})
	}
}

func TestCreateAPIJWT(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	varsEmptyAPIKey := *vars
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}
	varsEmptyAPIKey.APIKey = ""

	successfulCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should create a jwt api based",
			input:  "",
		},
	}

	failedCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should fail create a jwt api based, empty string",
			input:  "",
		},
		{
			config: &varsEmptyAPIKey,
			name:   "it should fail create a jwt api based, empty string",
			input:  "",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateAPIJWT()

			assert.NoError(t, err)
			assert.NotEmpty(t, r, "expected a non-empty string, but got empty")
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateUserJWT(tc.input)

			assert.Error(t, err)
			assert.Empty(t, r, "expected a non-empty string, but got empty")
		})
	}
}

func TestCheckAPIJWT(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	varsEmptyAPIKey := *vars
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}
	varsEmptyAPIKey.APIKey = ""

	successfulCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should check a jwt user based",
			input:  "im an id",
		},
	}

	failedCases := []struct {
		name   string
		input  string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should fail check a jwt user based, empty string",
			input:  "",
		},
		{
			config: &varsEmptyAPIKey,
			name:   "it should fail create a jwt api based, empty string",
			input:  "",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			r, err := jwt.CreateAPIJWT()

			if err != nil {
				t.Error("expected no error, but got:", err)
			}
			if r == "" {
				t.Error("expected a non-empty string, but got empty")
			}

			ok, err := jwt.CheckUserJwt(r)
			assert.NoError(t, err)
			assert.Equal(t, true, ok, "expected an ok, but got false")
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			jwt := utils.NewJWT(*tc.config)

			ok, err := jwt.CheckAPIJWT("")
			assert.Error(t, err)
			assert.Equal(t, false, ok, "expected an ok, but got false")
		})
	}
}
