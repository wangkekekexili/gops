package db

import (
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

// QuestionMarks generates something like (?,?,?) to be used by SQL statement.
// Caller must guarantee that input n is positive.
func QuestionMarks(n int) string {
	if n <= 0 {
		panic(fmt.Errorf("programming error: %v as input for QuestionMarks()", n))
	}
	return fmt.Sprintf("(%s?)", strings.Repeat("?,", n-1))
}
