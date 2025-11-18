package goda

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var pg *sql.DB

func GetPG(t *testing.T) *sql.DB {
	if pg == nil {
		t.Skip("pgx is not installed")
		t.SkipNow()
	}
	return pg
}

func init() {
	var e error
	pg, e = sql.Open("pgx", "")
	if e != nil {
		pg = nil
		return
	}
}
