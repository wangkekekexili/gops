package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB struct {
	Config *Config

	*sqlx.DB
}

func (db *DB) Load() error {
	if db.Config == nil || db.Config.MysqlDSN == "" {
		return errors.New("cannot load db")
	}
	d, err := sqlx.Connect("mysql", db.Config.MysqlDSN)
	if err != nil {
		return err
	}
	db.DB = d
	return nil
}
