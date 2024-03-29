package store

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *pgx.Conn
}

type People struct {
	ID   int
	Name string
}

// NewStore creates new database connection
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

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:../../migrations/1_initial.up.sql",
		"postgres", driver)
	if err != nil {
		panic(nil)
	}
	m.Up()

	return &Store{
		conn: conn,
	}
}

func (s *Store) ListPeople() ([]People, error) {
	rows, err := s.conn.Query(context.Background(), `
	SELECT * 
	FROM people
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]People, 5)
	for rows.Next() {
		var (
			id   string
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		i, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		res = append(res, People{
			ID:   i,
			Name: name,
		})

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return res, err
}

func (s *Store) GetPeopleByID(id string) (People, error) {
	var (
		pid   string
		pname string
	)
	err := s.conn.QueryRow(context.Background(), `
	SELECT id, name
	FROM people
	WHERE id = $1
	`, id).Scan(&pid, &pname)
	if err != nil {
		return People{}, err
	}
	i, err := strconv.Atoi(pid)
	if err != nil {
		return People{}, err
	}
	return People{
		ID:   i,
		Name: pname,
	}, nil
}
