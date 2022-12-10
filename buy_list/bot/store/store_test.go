package store

import (
	"buy_list/bot/models"
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

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

	connString := fmt.Sprintf("postgres://postgres:postgres@%v:%v/testdb?sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:../.././migrations/",
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
}

func TestGetUserByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetUserByUsername(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestCreatePrdouctByName(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	p1, err := store.CreateProductByName(context.Background(), "fish")
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}

	p2, err := store.CreateProductByName(context.Background(), "fish")
	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}

	if p1.ProductId != p2.ProductId {
		t.Error("а new product with the same name was created")
	}
}

func TestGetProductByName(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetProductByName(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestDeleteProdcutFromBuyListById(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.DeleteProductFromBuyListById(context.Background(), "b2a514a1-1755-4f9c-a3a3-132b5eb3a258", "b2a514a1-1755-4f9c-a3a3-132b5eb3a258")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestDeleteProductFromFridgeById(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.DeleteProductFromFridgeById(context.Background(), "b2a514a1-1755-4f9c-a3a3-132b5eb3a258", "b2a514a1-1755-4f9c-a3a3-132b5eb3a258")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestOpenProductFromFridgeById(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.DeleteProductFromFridgeById(context.Background(), "b2a514a1-1755-4f9c-a3a3-132b5eb3a258", "b2a514a1-1755-4f9c-a3a3-132b5eb3a258")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestSetCookedProductFromFridgeById(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.SetCookedProductFromFridgeById(context.Background(),
		"b2a514a1-1755-4f9c-a3a3-132b5eb3a258", time.Now().Format("2006-01-02"))

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestSetThrownProductFromFridgeById(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	err := store.SetThrownProductFromFridgeById(context.Background(),
		"b2a514a1-1755-4f9c-a3a3-132b5eb3a258", time.Now().Format("2006-01-02"))

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestAddProductToBuyList(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	p, _ := store.CreateProductByName(context.Background(), "fio")

	store.AddUsertg(context.Background(), &models.Usertg{
		Username: "aaaa",
	})

	u, _ := store.GetUserByUsername(context.Background(), "aaaa")

	err := store.AddProductToBuyList(context.Background(), &models.Product{
		UserId:    u.UserId,
		ProductId: p.ProductId,
		BuyDate:   "2006-01-02",
	})

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetBuyListByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetBuyListByUsername(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestAddProductToFridge(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	p, _ := store.CreateProductByName(context.Background(), "fio")

	store.AddUsertg(context.Background(), &models.Usertg{
		Username: "aaaa",
	})

	u, _ := store.GetUserByUsername(context.Background(), "aaaa")

	err := store.AddProductToFridge(context.Background(), &models.FridgeProduct{
		UserId:      u.UserId,
		ProductId:   p.ProductId,
		Expire_date: "2006-01-02",
	})

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetFridgeListByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetFridgeListByUsername(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetFridgeListByUsernameAlpha(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetFridgeListByUsernameAlpha(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetFridgeListByUsernameExpDate(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetFridgeListByUsernameExpDate(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetUsedProductsByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetUsedProductsByUsername(context.Background(), "aaaa")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetUsedProductsInPeriodByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetUsedProductsInPeriodByUsername(context.Background(), "aaaa",
		models.PeriodStat{
			FromDate: "2006-01-02",
			ToDate:   "2007-01-02",
		})

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetCountCookedUsedProductsInPeriodByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetCountCookedUsedProductsInPeriodByUsername(context.Background(), "aaaa",
		models.PeriodStat{
			FromDate: "2006-01-02",
			ToDate:   "2007-01-02",
		})

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetCountThrownUsedProductsInPeriodByUsername(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetCountThrownUsedProductsInPeriodByUsername(context.Background(), "aaaa",
		models.PeriodStat{
			FromDate: "2006-01-02",
			ToDate:   "2007-01-02",
		})

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetTodayBuyList(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetTodayBuyList(context.Background())

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetChatIdByUserId(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetChatIdByUserId(context.Background(), "b2a514a1-1755-4f9c-a3a3-132b5eb3a258")

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}

func TestGetSoonExpireList(t *testing.T) {
	container, d := CreateTestDatabase()
	defer container.Terminate(context.Background())

	store := Store{
		db: d,
	}

	_, err := store.GetSoonExpireList(context.Background())

	if err != nil && err != sql.ErrNoRows {
		t.Error(err)
	}
}
