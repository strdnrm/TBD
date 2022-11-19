package store

import (
	"buy_list/bot/models"
	"context"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateTestDatabase() (testcontainers.Container, *sqlx.DB) {
	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_USER":     "postgres",
		},
	}

	dbContainer, err := testcontainers.GenericContainer(
		context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: containerReq,
			Started:          true,
		})
	if err != nil {
		panic(err)
	}

	host, err := dbContainer.Host(context.Background())
	if err != nil {
		panic(err)
	}
	port, err := dbContainer.MappedPort(context.Background(), "5432")
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb", host, port.Port())
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations/",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()

	return dbContainer, db
}

func TestAddUsertg(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.AddUsertg(context.Background(), &models.Usertg{
		UserId:   "bbbbbbbb-1755-4f9c-a3a3-132b5eb3a258",
		Username: "aavasadas",
		ChatId:   1019622784,
	})

	if err != nil {
		t.Error(err)
	}
	// t.Run("AddUsertg", func(t *testing.T) {

	// })
}
