package goda

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlC *sql.DB

func GetMySQL(t *testing.T) *sql.DB {
	if mysqlC == nil {
		t.Skip("mysql is not installed")
		t.SkipNow()
	}
	return mysqlC
}

func init() {
	mysqlC, _ = sql.Open("mysql", "root:123456@/")
}
