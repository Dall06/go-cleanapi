// Package database_test is a test for db driver
package database_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/infrastructure/database"
	"dall06/go-cleanapi/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConn(test *testing.T) {
	cfg := config.NewConfig("8080", "0.0.0")
	vars, err := cfg.SetConfig()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	varsEmptyDBConn := *vars
	varsEmptyDBConn.DBConnString = ""

	logger := utils.NewLoggerMock()
	err = logger.Initialize()
	if err != nil {
		test.Fatal("expected no error, but got:", err)
	}

	successfulCases := []struct {
		name   string
		vars   *config.Vars
		logger utils.Logger
	}{
		{
			name:   "it should run and close a db conn",
			vars:   vars,
			logger: logger,
		},
	}

	failedCases := []struct {
		name   string
		vars   *config.Vars
		logger utils.Logger
	}{
		{
			name:   "it should not run a db conn, empty db conn string",
			vars:   &varsEmptyDBConn,
			logger: logger,
		},
	}

	failedCasesClose := []struct {
		name   string
		vars   config.Vars
		logger utils.Logger
	}{
		{
			name:   "it should not run a db conn, empty db conn string",
			vars:   *vars,
			logger: logger,
		},
	}

	for _, tc := range successfulCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := database.NewDBConn(tc.logger, *tc.vars)
			conn, err := db.Open()
			assert.NoError(t, err)
			assert.NotEmpty(t, conn, "connection should not be empty")

			err = db.Close(conn)
			assert.NoError(t, err)
		})
	}

	for _, tc := range failedCases {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := database.NewDBConn(tc.logger, *tc.vars)
			conn, err := db.Open()
			assert.NotEmpty(t, err, "expected error, but got nil")
			assert.Empty(t, conn, "connection should be empty")

			err = db.Close(conn)
			assert.Error(t, err)
		})
	}

	for _, tc := range failedCasesClose {
		tc := tc

		test.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := database.NewDBConn(tc.logger, tc.vars)

			err = db.Close(nil)
			assert.Error(t, err)
		})
	}
}
