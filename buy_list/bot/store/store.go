package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	conn *pgx.Conn
	db   *sqlx.DB
}
type Product struct { // 0 - name ; 1 - weight ; 2 - buydate
	UserId    string
	ProductId string `db:"id"`
	State     int
	Name      string `db:"name"`
	Weight    float64
	BuyDate   string
}

type FridgeProduct struct { // 0 - name ; 1 - expire date
	UserId      string `db:"user_id"`
	ProductId   string `db:"product_id"`
	State       int
	Name        string
	Opened      bool   `db:"opened"`
	Expire_date string `db:"expire_date"`
	Status      string `db:"status"`
	Use_date    string `db:"use_date"`
}

type Usertg struct {
	UserId   string `db:"id"`
	Username string `db:"username"`
}

type BuyList struct {
	UserId    string  `db:"user_id"`
	ProductId string  `db:"product_id"`
	Weight    float64 `db:"weight"`
	BuyTime   string  `db:"buy_time"`
}

//TODO inteface

func NewStore(connString string) *Store {

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return &Store{
		conn: conn,
		db:   db,
	}
}

func (s *Store) AddUsertg(ctx context.Context, u *Usertg) error {
	rows, err := s.conn.Query(ctx, `
	INSERT INTO usertg(username)
	VALUES ($1);
	`, u.Username)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (s *Store) GetUseridByUsername(ctx context.Context, username string) (string, error) {
	var id string
	err := s.conn.QueryRow(ctx, `
	SELECT id::text FROM usertg WHERE username = $1;
	`, username).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Store) CreateProductByName(ctx context.Context, productName string) (string, error) {
	var id string

	err := s.conn.QueryRow(ctx, `
	INSERT INTO product(name)
	VALUES($1) RETURNING id::text;
	`, productName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Store) GetProductIdByName(ctx context.Context, productName string) (string, error) {
	var id string
	err := s.conn.QueryRow(ctx, `
	SELECT id FROM product
	WHERE name = $1
	`, productName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Store) DeleteProductFromBuyListById(ctx context.Context, productId string) error {
	rows, err := s.conn.Query(ctx, `
	DELETE FROM buy_list
	WHERE product_id = $1
	`, productId)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (s *Store) DeleteProductFromFridgeById(ctx context.Context, productId string) error {
	rows, err := s.conn.Query(ctx, `
	DELETE FROM fridge
	WHERE product_id = $1
	`, productId)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (s *Store) OpenProductFromFridgeById(ctx context.Context, productId string, expDate string) error {
	rows, err := s.conn.Query(ctx, `
	UPDATE fridge 
	SET opened = true, expire_date = $1
	WHERE product_id = $2
	`, expDate, productId)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (s *Store) SetCookedProductFromFridgeById(ctx context.Context, productId string, useDate string) error {
	rows, err := s.conn.Query(ctx, `
	UPDATE fridge 
	SET status = 'cooked', use_date = $1
	WHERE product_id = $2
	`, useDate, productId)
	if err != nil {
		return err
	}
	rows.Close()
	return nil
}

func (s *Store) SetThrownProductFromFridgeById(ctx context.Context, productId string, useDate string) error {
	rows, err := s.conn.Query(ctx, `
	UPDATE fridge 
	SET status = 'thrown', use_date = $1
	WHERE product_id = $2
	`, useDate, productId)
	rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddProductToBuyList(ctx context.Context, p *Product) error {
	rows, err := s.conn.Query(ctx, `
	INSERT INTO buy_list
	VALUES($1, $2, $3, $4)
	`, p.UserId, p.ProductId, p.Weight, p.BuyDate)
	rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBuyListByUsername(ctx context.Context, username string) ([]Product, error) {
	// get name wight buydate
	rows, err := s.conn.Query(ctx, `
	SELECT product.name, buy_list.weight, buy_list.buy_time FROM buy_list
	JOIN product ON product.id = buy_list.product_id
	JOIN usertg ON usertg.id = buy_list.user_id
	WHERE usertg.username = $1;
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Product
	if rows.Err() == pgx.ErrNoRows {
		return list, nil
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
	return list, nil
}

func (s *Store) AddProductToFridge(ctx context.Context, f *FridgeProduct) error {
	rows, err := s.conn.Query(ctx, `
	INSERT INTO fridge
	VALUES($1, $2,
	FALSE, $3, NULL, NULL)
	`, f.UserId, f.ProductId, f.Expire_date)
	rows.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetFridgeListByUsername(ctx context.Context, username string) ([]FridgeProduct, error) {
	//get name opened expire_date  status
	rows, err := s.conn.Query(ctx, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FridgeProduct
	if rows.Err() == pgx.ErrNoRows {
		return list, nil
	}
	for rows.Next() {
		f := FridgeProduct{}
		// var tmppp pgtype.Timestamptz
		var expDate time.Time
		if err := rows.Scan(&f.Name, &f.Opened, &expDate); err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println("Failed scan ", err)
		}
		f.Expire_date = expDate.Format("2006-01-02")
		// f.Use_date = useDate.Format("2006-02-01")
		list = append(list, f)

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return list, nil
}

func (s *Store) GetFridgeListByUsernameAlpha(ctx context.Context, username string) ([]FridgeProduct, error) {
	//get name opened expire_date  status
	rows, err := s.conn.Query(ctx, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	ORDER BY pd.name
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FridgeProduct
	if rows.Err() == pgx.ErrNoRows {
		return list, nil
	}
	for rows.Next() {
		f := FridgeProduct{}
		// var tmppp pgtype.Timestamptz
		var expDate time.Time
		if err := rows.Scan(&f.Name, &f.Opened, &expDate); err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println("Failed scan ", err)
		}
		f.Expire_date = expDate.Format("2006-01-02")
		// f.Use_date = useDate.Format("2006-02-01")
		list = append(list, f)

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return list, nil
}

func (s *Store) GetFridgeListByUsernameExpDate(ctx context.Context, username string) ([]FridgeProduct, error) {
	//get name opened expire_date  status
	rows, err := s.conn.Query(ctx, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	ORDER BY f.expire_date, pd.name
	`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FridgeProduct
	if rows.Err() == pgx.ErrNoRows {
		return list, nil
	}
	for rows.Next() {
		f := FridgeProduct{}
		// var tmppp pgtype.Timestamptz
		var expDate time.Time
		if err := rows.Scan(&f.Name, &f.Opened, &expDate); err != nil {
			var Error = log.New(os.Stdout, "\u001b[31mERROR: \u001b[0m", log.LstdFlags|log.Lshortfile)
			Error.Println("Failed scan ", err)
		}
		f.Expire_date = expDate.Format("2006-01-02")
		// f.Use_date = useDate.Format("2006-02-01")
		list = append(list, f)

		if rows.Err() != nil {
			fmt.Fprintf(os.Stderr, "Scan error: %v\n", rows.Err())
		}
	}
	return list, nil
}
