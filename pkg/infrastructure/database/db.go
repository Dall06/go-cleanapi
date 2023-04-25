// Package database is the package that contains the sql driver
package database

import (
	"dall06/go-cleanapi/config"
	"dall06/go-cleanapi/utils"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql" // this package registers the mysql driver with sql
)

// DB is an interface that extend dbConn
type DB interface {
	Open() (*sql.DB, error)
	Close(db *sql.DB) error
}

const (
	emptyConnectionString = "empty connection string"
	emptydbConn           = "empty connection"
	dbEngine              = "mysql"
	maxLifeTime           = time.Minute * 3
	maxOpenConns          = 10
	idleConns             = 10
)

var _ DB = (*dbConn)(nil)

type dbConn struct {
	logger utils.Logger
	config config.Vars
}

// NewDBConn is a constructor for dbConn
func NewDBConn(l utils.Logger, v config.Vars) DB {
	return &dbConn{
		logger: l,
		config: v,
	}
}

func (c *dbConn) Open() (*sql.DB, error) {
	if c.config.DBConnString == "" {
		c.logger.Warn("db connection failed: %v", emptyConnectionString)
		return nil, errors.New(emptyConnectionString)
	}

	db, err := sql.Open(dbEngine, c.config.DBConnString)
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

func (c *dbConn) Close(db *sql.DB) error {
	if db == nil {
		c.logger.Warn("db connection close failed: %s", emptyConnectionString)
		return errors.New(emptydbConn)
	}

	err := db.Close()
	if err != nil {
		c.logger.Error("db connection close failed: %v", err)
		return err
	}

	//c.logger.Info("db connection closed", nil)
	return nil
}
