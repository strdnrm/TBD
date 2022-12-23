package sqlstore

import (
	"context"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestDB(t *testing.T, ctx context.Context, databaseURL string) (*sqlx.DB, func(...string)) {
	t.Helper()

	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.ExecContext(ctx, "TRUNCATE %s CASCADE", strings.Join(tables, ", "))
		}

		db.Close()
	}
}
