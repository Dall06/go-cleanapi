package utils_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"fmt"
	"os"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestLoggerDev(test *testing.T) {
	conf := config.NewConfig("8080")
	err := conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to create config: %v", err)
	}

	l := utils.NewLogger()
	if l == nil {
		test.Fatalf("failed to load logger")
	}

	err = l.Initialize()
	if err != nil {
		test.Fatalf("failed to initialize logger: %v", err)
	}

	test.Run("it should log in warn", func(t *testing.T) {
		// Create a temporary log file
		msg := "test warning message"
		l.Warn(msg)

		filePath := fmt.Sprintf("../logs/%s_%s.log", config.Stage, zap.WarnLevel.String())
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("error reading log file to check content: %v", err)
		}

		if !strings.Contains(string(data), msg) {
			t.Errorf("log file doesn't contain expected message: %q", msg)
		}
	})

	test.Run("it should log in error", func(t *testing.T) {
		// Create a temporary log file
		msg := "test warning message"
		l.Error(msg)
		l.Error("GET %s:%s", "internalError", "err")

		filePath := fmt.Sprintf("../logs/%s_%s.log", config.Stage, zap.WarnLevel.String())
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("error reading log file to check content: %v", err)
		}

		if !strings.Contains(string(data), msg) {
			t.Errorf("log file doesn't contain expected message: %q", msg)
		}
	})

	test.Run("it should log in info", func(t *testing.T) {
		// Create a temporary log file
		msg := "test warning message"
		l.Info(msg)
		l.Info("GET %s:%s", "im a get", "accessed")

		filePath := fmt.Sprintf("../logs/%s_%s.log", config.Stage, zap.WarnLevel.String())
		data, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("error reading log file to check content: %v", err)
		}

		if !strings.Contains(string(data), msg) {
			t.Errorf("log file doesn't contain expected message: %q", msg)
		}
	})

	test.Run("it should not log, wrong dir", func(t *testing.T) {
		// Create a temporary log file
		msg := "test warning message"
		l.Warn(msg)
		l.Warn("GET %s:%s", "im a get", "something might be bad")

		filePath := fmt.Sprintf("/logs/%s_%s.log", config.Stage, zap.WarnLevel.String())
		_, err := os.ReadFile(filePath)
		if err == nil {
			t.Fatalf("i must have failed")
		}
	})
}
