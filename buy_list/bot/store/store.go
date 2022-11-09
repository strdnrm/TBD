package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *pgx.Conn
}

type ApiToken struct {
	Tkn string `json:"token"`
}

type DBdata struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Product struct { // 0 - name ; 1 - weight ; 2 - buydate
	UserId    string
	ProductId string
	State     int
	Name      string
	Weight    float64
	BuyDate   string
}

type FridgeProduct struct { // 0 - name ; 1 - expire date
	UserId      string
	ProductId   string
	State       int
	Name        string
	Opened      bool
	Expire_date string
	Status      string
	Use_date    string
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

func (s *Store) GetUseridByUsername(username string) string {
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

func (s *Store) CreateProductByName(productName string) string {
	var id string
	err := s.conn.QueryRow(context.Background(), `
	INSERT INTO product(name)
	VALUES($1) RETURNING id::text;
	`, productName).Scan(&id)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Get product id error while creating", err)
	}

	return id
}

func (s *Store) GetProductIdByName(productName string) string {
	var id string
	err := s.conn.QueryRow(context.Background(), `
	SELECT id FROM product
	WHERE name = $1
	`, productName).Scan(&id)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Get product id error ", err)
	}

	return id
}

func (s *Store) DeleteProductFromBuyListById(productId string) {
	rows, err := s.conn.Query(context.Background(), `
	DELETE FROM buy_list
	WHERE product_id = $1
	`, productId)
	rows.Close()
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Delete product error ", err)
	}
}

func (s *Store) AddProductToBuyList(p *Product) {
	rows, err := s.conn.Query(context.Background(), `
	INSERT INTO buy_list
	VALUES($1, $2, $3, $4)
	`, p.UserId, p.ProductId, p.Weight, p.BuyDate)
	rows.Close()
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Add product to buy list error ", err)
	}

}

func (s *Store) GetBuyListByUsername(username string) []Product {
	// get name wight buydate
	rows, err := s.conn.Query(context.Background(), `
	SELECT product.name, buy_list.weight, buy_list.buy_time FROM buy_list
	JOIN product ON product.id = buy_list.product_id
	JOIN usertg ON usertg.id = buy_list.user_id
	WHERE usertg.username = $1;
	`, username)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Getting buy list ", err)
	}
	defer rows.Close()
	var list []Product
	if rows.Err() == pgx.ErrNoRows {
		return list
	}
	for rows.Next() {
		p := Product{}
		// var tmppp pgtype.Timestamptz
		var tmpTime time.Time
		if err := rows.Scan(&p.Name, &p.Weight, &tmpTime); err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println("Failed scan ", err)
		}
		p.BuyDate = tmpTime.Format(time.RFC850)
		list = append(list, p)

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return list
}

func (s *Store) AddProductToFridge(f *FridgeProduct) {
	rows, err := s.conn.Query(context.Background(), `
	INSERT INTO fridge
	VALUES($1, $2,
	FALSE, $3, NULL, NULL)
	`, f.UserId, f.ProductId, f.Expire_date)
	rows.Close()
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Add product to fridge error ", err)
	}
}

func (s *Store) GetFridgeListByUsername(username string) []FridgeProduct {
	//get name opened expire_date  status
	rows, err := s.conn.Query(context.Background(), `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	`, username)
	if err != nil {
		var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
		Error.Println("Getting fridge list ", err)
	}
	defer rows.Close()
	var list []FridgeProduct
	if rows.Err() == pgx.ErrNoRows {
		return list
	}
	for rows.Next() {
		f := FridgeProduct{}
		// var tmppp pgtype.Timestamptz
		var expDate time.Time
		if err := rows.Scan(&f.Name, &f.Opened, &expDate); err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println("Failed scan ", err)
		}
		fmt.Println("ALLEEE!!!" + f.Name)
		f.Expire_date = expDate.Format("2006-02-01")
		// f.Use_date = useDate.Format("2006-02-01")
		list = append(list, f)

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return list
}
