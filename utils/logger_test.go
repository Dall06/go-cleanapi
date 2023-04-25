// Package utils_test is a test package for utils
package utils_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeLogger(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	varsEmptyStage := *vars
	varsEmptyStage.Stage = ""

	successfulCases := []struct {
		name   string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should initialize logger",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := utils.NewLogger(*tc.config)
			err := logger.Initialize()
			assert.NoError(t, err)
		})
	}
}

func TestLoggerWarn(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should log a warn",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := utils.NewLogger(*tc.config)
			err := logger.Initialize()
			assert.NoError(t, err)

			logger.Warn("im a format %s", "im a string")
		})
	}
}

func TestLoggerInfo(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should log info",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := utils.NewLogger(*tc.config)
			err := logger.Initialize()
			assert.NoError(t, err)

			logger.Info("im a format %s", "im a string")
		})
	}
}

func TestLoggerError(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		config *config.Vars
	}{
		{
			config: vars,
			name:   "it should log a error",
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := utils.NewLogger(*tc.config)
			err := logger.Initialize()
			assert.NoError(t, err)

			logger.Error("im a format %s", "im a string")
		})
	}
}
