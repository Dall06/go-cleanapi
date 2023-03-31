package database

import (
	"dall06/go-cleanapi/config"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBRepository interface {
	Open() (sql.DB, error)
	Close(db sql.DB) error
}

var _ DBRepository = (*DBConn)(nil)

type DBConn struct{}

func NewDBConn() DBConn {
	return DBConn{}
}

func (DBConn) Open() (sql.DB, error) {
	db, err := sql.Open("mysql", config.DBConnString)
	if err != nil {
		return *db, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return *db, nil
}

func (DBConn) Close(db sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
