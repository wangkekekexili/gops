package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wangkekekexili/gops/server"
	"github.com/wangkekekexili/module"
)

func main() {
	s := &server.GOPS{}
	err := module.Load(s)
	if err != nil {
		fmt.Printf("load failed: %v\n", err)
		os.Exit(1)
	}
	err = s.Start()
	if err != nil {
		s.Reporter.ErrSMS(err)
	}
}
