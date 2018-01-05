package db

import (
	"errors"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/wangkekekexili/gops/till"
)

type Module struct {
	Till *till.Module

	*sqlx.DB
}

func (m *Module) Load() (err error) {
	defer func() {
		if err != nil {
			m.Till.Notify(err)
		}
	}()

	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		return errors.New("MySQL_DSN is empty")
	}
	d, err := sqlx.Connect("mysql", mysqlDSN)
	if err != nil {
		return err
	}
	m.DB = d
	return nil
}
