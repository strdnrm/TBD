package apiserver

import (
	"final_project/internal/app/store/sqlstore"
	"net/http"

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

	// defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
