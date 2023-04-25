// Package config_test is a package to test the config variables
package config_test

import (
	"dall06/go-cleanapi/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigVars(test *testing.T) {
	successfulCases := []struct {
		name    string
		port    string
		version string
	}{
		{
			name:    "it should load all variables",
			port:    "8080",
			version: "0",
		},
	}

	failedCases := []struct {
		name    string
		port    string
		version string
	}{
		{
			name:    "it should not load all variables, empty port",
			port:    "",
			version: "0",
		},
		{
			name:    "it should not load all variables, empty version",
			port:    "8080",
			version: "",
		},
	}

	for _, tc := range successfulCases {
		tc := tc
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.NewConfig(tc.port, tc.version)
			vars, err := cfg.SetConfig()
			assert.NoError(t, err)

			assert.NotEmpty(t, vars, "expected vars, but got nil")

			assert.NotEmpty(t, vars.APIBasePath, "expected api base path, but got empty")
			assert.NotEmpty(t, vars.APIPort, "expected api port, but got empty")
			assert.NotEmpty(t, vars.APIKey, "expected api key, but got empty")
			assert.NotEmpty(t, vars.APIKeyHash, "expected api hash key, but got empty")
			assert.NotEmpty(t, vars.DBConnString, "expected db conn str, but got empty")
			assert.NotEmpty(t, vars.JWTSecret, "expected jwt secret, but got empty")
			assert.NotEmpty(t, vars.ProyectName, "expected proyect name, but got empty")
			assert.NotEmpty(t, vars.Stage, "expected stage, but got empty")
			assert.NotEmpty(t, vars.ProyectPath, "expected proyect path, but got empty")
			assert.NotEmpty(t, vars.CookieSecret, "expected cookie secret, but got empty")
			assert.NotEmpty(t, vars.APIVersion, "expected api version, but got empty")
			assert.NotEmpty(t, vars.AppName, "expected app name, but got empty")
		})
	}

	for _, tc := range failedCases {
		tc := tc
		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.NewConfig(tc.port, tc.version)
			vars, err := cfg.SetConfig()
			assert.Error(t, err)
			assert.Empty(t, vars, "expected nil, but got vars")
		})
	}
}
