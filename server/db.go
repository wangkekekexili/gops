package server

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type DB struct {
	Config *Config

	*sql.DB
}

func (db *DB) Load() error {
	if db.Config == nil || db.Config.MysqlDSN == "" {
		return errors.New("cannot load db")
	}
	d, err := sql.Open("mysql", db.Config.MysqlDSN)
	if err != nil {
		return err
	}
	db.DB = d

	err = d.Ping()
	if err != nil {
		return err
	}
	return nil
}
