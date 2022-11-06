package store

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *pgx.Conn
}

type Product struct { // 0 - name ; 1 - weight ; 2 - buydate
	UserId    string
	ProductId string
	State     int
	Name      string
	Weight    float64
	BuyDate   string
}

//TODO inteface

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

func (s *Store) AddUsertg(u *Usertg) {
	rows, err := s.conn.Query(context.Background(), `
	INSERT INTO usertg(username)
	VALUES ($1);
	`, u.Username)
	defer rows.Close()
	if err != nil {
		//TODO check error types
		//var pgErr

		Warning := log.New(os.Stdout, "\u001b[33mWARNING: \u001B[0m", log.LstdFlags|log.Lshortfile)
		Warning.Println("Username already exists")
	}
}

func (s *Store) GetUserid(username string) string {
	var id string
	err := s.conn.QueryRow(context.Background(), `
	SELECT id::text FROM usertg WHERE username = $1;
	`, username).Scan(&id)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Get userid error ", err)
	}

	return id
}

func (s *Store) GetProductId(productName string) string {
	var id string
	err := s.conn.QueryRow(context.Background(), `
	INSERT INTO product(name)
	VALUES($1) RETURNING id::text;
	`, productName).Scan(&id)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Get product id error ", err)
	}

	return id
}

func (s *Store) AddProductToBuyList(p *Product) {
	rows, err := s.conn.Query(context.Background(), `
	INSERT INTO buy_list
	VALUES($1, $2, $3, $4)
	`, p.UserId, p.ProductId, p.Weight, p.BuyDate)
	rows.Close()
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Add product error ", err)
	}

}
