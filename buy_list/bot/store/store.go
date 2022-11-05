package store

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *pgx.Conn
}

type Product struct {
	//	UUID string
	Name string
}

type Usertg struct {
	//	UUID     string
	Username string
}

func NewStore(connString string) *Store {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return &Store{
		conn: conn,
	}
}

func (s *Store) AddUsertg(u Usertg) {
	err := s.conn.QueryRow(context.Background(), `
	INSERT INTO usertg(username)
	VALUES ($1);
	`, u.Username)
	if err != nil {
		log.Panic(err)
	}
}
