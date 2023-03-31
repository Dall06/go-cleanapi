package database_test

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/pkg/infrastructure/database"
	"testing"
)

func TestDatabase(test *testing.T) {
	conf := config.NewConfig()
	err := conf.SetConfig()
	if err != nil {
		test.Fatalf("failed to load: %s", err)
	}

	conn := database.NewDBConn()
	test.Run("it should connect to test db", func(t *testing.T) {
		_, err := conn.Open()
		if err != nil {
			test.Fatalf("open db failed: %v", err)
		}
	})

	test.Run("it should connect and close connection to test db", func(t *testing.T) {
		db, err := conn.Open()
		if err != nil {
			test.Fatalf("open db failed: %v", err)
		}

		err = conn.Close(db)
		if err != nil {
			test.Fatalf("close db failed: %v", err)
		}
	})
}
