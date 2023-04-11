package database

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
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

type DBConn struct{
	logger utils.Logger
}

func NewDBConn(l utils.Logger) DBRepository {
	return DBConn{
		logger: l,
	}
}

func (c DBConn) Open() (*sql.DB, error) {
	if config.DBConnString == "" {
		c.logger.Warn("db connection failed: %v", emptyConnectionString)
		return nil, errors.New(emptyConnectionString)
	}

	db, err := sql.Open(dbEngine, config.DBConnString)
	if err != nil {
		c.logger.Error("db connection failed: %v", err)
		return nil, err
	}

	db.SetConnMaxLifetime(maxLifeTime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(idleConns)

	c.logger.Info("db connection opened")
	return db, nil
}

func (c DBConn) Close(db *sql.DB) error {
	if db == nil {
		c.logger.Warn("db connection close failed: %s", emptyConnectionString)
		return errors.New(emptyDBConn)
	}

	err := db.Close()
	if err != nil {
		c.logger.Error("db connection close failed: %v", err)
		return err
	}

	c.logger.Info("db connection closed")
	return nil
}
