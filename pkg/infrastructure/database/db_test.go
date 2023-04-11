package database_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/infrastructure/database"
	"dall06/go-cleanapi/utils"
	"testing"
)

func TestDatabase(test *testing.T) {
	lm := utils.NewMockLoggerRepository()
	lm.Initialize()

	test.Run("it should not connect to test db", func(t *testing.T) {
		t.Parallel()
		conn := database.NewDBConn(lm)

		_, err := conn.Open()
		if err == nil {
			test.Fatalf("open db failed: %v", err)
		}
	})

	test.Run("it should connect and close connection to test db", func(t *testing.T) {
		t.Parallel()
		conf := config.NewConfig("8080")
		err := conf.SetConfig()
		if err != nil {
			test.Fatalf("failed to load: %s", err)
		}

		conn := database.NewDBConn(lm)

		db, err := conn.Open()
		if err != nil {
			test.Fatalf("open db failed: %v", err)
		}

		err = conn.Close(db)
		if err != nil {
			test.Fatalf("close db failed: %v", err)
		}
	})

	test.Run("it should not connect and close connection to test db, bc empty db", func(t *testing.T) {
		t.Parallel()
		conn := database.NewDBConn(lm)

		err := conn.Close(nil)
		if err == nil {
			test.Fatalf("close db failed: %v", err)
		}
	})
}
