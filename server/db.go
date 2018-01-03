package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DB struct {
	Config *Config
	Till   *Till

	*sqlx.DB
}

func (db *DB) Load() (err error) {
	defer func() {
		if err != nil {
			db.Till.Notify(err)
		}
	}()

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
