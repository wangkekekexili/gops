package server_test

import (
	"os"
	"testing"

	"github.com/wangkekekexili/gops/server"
)

func TestConfig(t *testing.T) {
	c := &server.Config{}
	err := c.Load()
	if err == nil {
		t.Fatal("expect error when required values are not present")
	}

	os.Setenv("MYSQL_DSN", "mysql")
	os.Setenv("SENTRY_DSN", "sentry")
	c = &server.Config{}
	err = c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if c.MysqlDSN != "mysql" {
		t.Fatalf("unexpected MYSQL_DSN %v", c.MysqlDSN)
	}
	if c.SentryDSN != "sentry" {
		t.Fatalf("unexpected SENTRY_DSN %v", c.SentryDSN)
	}
}
