package utils_test

import (
	"dall06/go-cleanapi/utils"
	"testing"
)

func TestLoggerDev(test *testing.T) {
	l := utils.NewLogger("dev")
	if l == nil {
		test.Fatalf("failed to load logger")
	}
	test.Run("it should log in warn", func(t *testing.T) {
		err := l.Warn("test warn: ", "hi")
		if err != nil {
			t.Fatalf("failed to log warn: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in info", func(t *testing.T) {
		err := l.Info("test info: ", "hi")
		if err != nil {
			t.Fatalf("failed to log info: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in error", func(t *testing.T) {
		err := l.Error("test error: ", "hi")
		if err != nil {
			t.Fatalf("failed to log error: %s", err)
		}

		t.Log("log success ", "hi")
	})
}

func TestLoggerProd(test *testing.T) {
	l := utils.NewLogger("prod")
	if l == nil {
		test.Fatalf("failed to load logger")
	}
	test.Run("it should log in warn", func(t *testing.T) {
		err := l.Warn("test warn: ", "hi")
		if err != nil {
			t.Fatalf("failed to log warn: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in info", func(t *testing.T) {
		err := l.Info("test info: ", "hi")
		if err != nil {
			t.Fatalf("failed to log info: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in error", func(t *testing.T) {
		err := l.Error("test error: ", "hi")
		if err != nil {
			t.Fatalf("failed to log error: %s", err)
		}

		t.Log("log success  ", "hi")
	})
}

func TestLoggerTest(test *testing.T) {
	l := utils.NewLogger("test")
	if l == nil {
		test.Fatalf("failed to load logger")
	}
	test.Run("it should log in warn", func(t *testing.T) {
		err := l.Warn("test warn: ", "hi")
		if err != nil {
			t.Fatalf("failed to log warn: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in info", func(t *testing.T) {
		err := l.Info("test info: ", "hi")
		if err != nil {
			t.Fatalf("failed to log info: %s", err)
		}

		t.Log("log success ", "hi")
	})

	test.Run("it should log in error", func(t *testing.T) {
		err := l.Error("test error: ", "hi")
		if err != nil {
			t.Fatalf("failed to log error: %s", err)
		}

		t.Log("log failed succesfully ", err)
	})
}

func TestLog(test *testing.T) {
	test.Run("it should return an empty logger error", func(t *testing.T) {
		l := utils.NewLogger("non_exist")
		if l != nil {
			test.Fatalf("failed to fail")
		}

		test.Log("log failed succesfully: empty logger")
	})
}
