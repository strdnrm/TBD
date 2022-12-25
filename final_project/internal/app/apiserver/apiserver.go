package apiserver

import (
	"final_project/internal/app/store/sqlstore"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))

	store := sqlstore.New(db)
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.Addr, srv)
}

func newDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations/",
		"postgres", driver,
	)
	if err != nil {
		return nil, err
	}

	m.Up()

	return db, nil
}
