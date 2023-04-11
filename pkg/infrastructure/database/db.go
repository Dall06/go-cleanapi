package database

import (
	"dall06/go-cleanapi/config"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBRepository interface {
	Open() (*sql.DB, error)
	Close(db *sql.DB) error
}

const (
	emptyConnectionString = "empty connection string"
	emptyDBConn = "empty connection"
	dbEngine = "mysql"
	maxLifeTime = time.Minute * 3
	maxOpenConns = 10
	idleConns = 10
)

var _ DBRepository = (*DBConn)(nil)

type DBConn struct{}

func NewDBConn() DBRepository {
	return DBConn{}
}

func (DBConn) Open() (*sql.DB, error) {
	if config.DBConnString == "" {
		return nil, errors.New(emptyConnectionString)
	}

	db, err := sql.Open(dbEngine, config.DBConnString)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(maxLifeTime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(idleConns)

	return db, nil
}

func (DBConn) Close(db *sql.DB) error {
	if db == nil {
		return errors.New(emptyDBConn)
	}

	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
