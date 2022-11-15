package store

import (
	"context"

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
	UserId    string `db:"user_id"`
	ProductId string `db:"id"`
	State     int
	Name      string  `db:"name"`
	Weight    float64 `db:"weight"`
	BuyDate   string  `db:"buy_time"`
}

type FridgeProduct struct { // 0 - name ; 1 - expire date
	UserId      string `db:"user_id"`
	ProductId   string `db:"product_id"`
	State       int
	Name        string `db:"name"`
	Opened      bool   `db:"opened"`
	Expire_date string `db:"expire_date"`
	Status      string `db:"status"`
	Use_date    string `db:"use_date"`
}

type Usertg struct {
	UserId   string `db:"id"`
	Username string `db:"username"`
	ChatId   int64  `db:"chat_id"`
}

type PeriodStat struct {
	State    int
	FromDate string
	ToDate   string
}

// type BuyListProduct struct {
// 	UserId    string  `db:"user_id"`
// 	ProductId string  `db:"product_id"`
// 	Weight    float64 `db:"weight"`
// 	BuyTime   string  `db:"buy_time"`
// 	Name      string  `db:"name"`
// }

//TODO inteface

func NewStore(connString string) *Store {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	return &Store{
		conn: conn,
		db:   db,
	}
}

func (s *Store) AddUsertg(ctx context.Context, u *Usertg) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO usertg(username, chat_id)
	VALUES ($1, $2);
	`, u.Username, u.ChatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (Usertg, error) {
	u := Usertg{}
	err := s.db.GetContext(ctx, &u, `
	SELECT id::text FROM usertg WHERE username = $1;
	`, username)
	if err != nil {
		return u, err
	}
	return u, nil
}

// problem with the same product name
func (s *Store) CreateProductByName(ctx context.Context, productName string) (Product, error) {
	p := Product{}

	err := s.db.GetContext(ctx, &p, `
	INSERT INTO product(name)
	VALUES($1) RETURNING id::text;
	`, productName)
	p.Name = productName
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *Store) GetProductByName(ctx context.Context, productName string) (Product, error) {
	p := Product{}
	err := s.db.GetContext(ctx, &p, `
	SELECT id FROM product
	WHERE name = $1;
	`, productName)
	if err != nil {
		return p, err
	}
	return p, nil

}

func (s *Store) DeleteProductFromBuyListById(ctx context.Context, productId string) error {
	_, err := s.db.ExecContext(ctx, `
	DELETE FROM buy_list
	WHERE product_id = $1;
	`, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteProductFromFridgeById(ctx context.Context, productId string) error {
	_, err := s.db.ExecContext(ctx, `
	DELETE FROM fridge
	WHERE product_id = $1;
	`, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) OpenProductFromFridgeById(ctx context.Context, productId string, expDate string) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE fridge 
	SET opened = true, expire_date = $1
	WHERE product_id = $2;
	`, expDate, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SetCookedProductFromFridgeById(ctx context.Context, productId string, useDate string) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE fridge 
	SET status = 'cooked', use_date = $1
	WHERE product_id = $2;
	`, useDate, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SetThrownProductFromFridgeById(ctx context.Context, productId string, useDate string) error {
	_, err := s.db.ExecContext(ctx, `
	UPDATE fridge 
	SET status = 'thrown', use_date = $1
	WHERE product_id = $2;
	`, useDate, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddProductToBuyList(ctx context.Context, p *Product) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO buy_list
	VALUES($1, $2, $3, $4);
	`, p.UserId, p.ProductId, p.Weight, p.BuyDate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBuyListByUsername(ctx context.Context, username string) ([]Product, error) {
	// get name wight buydate
	var list []Product
	err := s.db.SelectContext(ctx, &list, `
	SELECT product.name, buy_list.weight, buy_list.buy_time FROM buy_list
	JOIN product ON product.id = buy_list.product_id
	JOIN usertg ON usertg.id = buy_list.user_id
	WHERE usertg.username = $1;
	`, username)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) AddProductToFridge(ctx context.Context, f *FridgeProduct) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO fridge
	VALUES($1, $2,
	FALSE, $3, NULL, NULL);
	`, f.UserId, f.ProductId, f.Expire_date)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetFridgeListByUsername(ctx context.Context, username string) ([]FridgeProduct, error) {
	var list []FridgeProduct
	err := s.db.SelectContext(ctx, &list, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null;
	`, username)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) GetFridgeListByUsernameAlpha(ctx context.Context, username string) ([]FridgeProduct, error) {
	var list []FridgeProduct
	err := s.db.SelectContext(ctx, &list, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	ORDER BY pd.name;
	`, username)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) GetFridgeListByUsernameExpDate(ctx context.Context, username string) ([]FridgeProduct, error) {
	//get name opened expire_date  status
	var list []FridgeProduct
	err := s.db.SelectContext(ctx, &list, `
	SELECT pd.name, f.opened, f.expire_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS null
	ORDER BY f.expire_date, pd.name;
	`, username)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (s *Store) GetUsedProductsByUsername(ctx context.Context, username string) ([]FridgeProduct, error) {
	products := []FridgeProduct{}
	err := s.db.SelectContext(ctx, &products, `
	SELECT pd.name, f.status, f.use_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1 AND f.status IS NOT null
	ORDER BY f.use_date;
	`, username)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *Store) GetUsedProductsInPeriodByUsername(ctx context.Context,
	username string, fromDate string, toDate string) ([]FridgeProduct, error) {
	products := []FridgeProduct{}
	err := s.db.SelectContext(ctx, &products, `
	SELECT pd.name, f.status, f.use_date
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	JOIN product pd ON pd.id = f.product_id
	WHERE ut.username = $1
		AND f.status IS NOT null
		AND f.use_date >= $2
		AND f.use_date <= $3
	ORDER BY f.use_date;
	`, username, fromDate, toDate)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *Store) GetCountCookedUsedProductsInPeriodByUsername(ctx context.Context,
	username string, fromDate string, toDate string) (int, error) {
	var cookedCount int
	err := s.db.GetContext(ctx, &cookedCount, `
	SELECT COUNT(*)
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	WHERE ut.username = $1
	AND f.status = 'cooked'
		AND f.use_date >= $2
		AND f.use_date <= $3;
	`, username, fromDate, toDate)
	if err != nil {
		return -1, err
	}
	return cookedCount, nil
}

func (s *Store) GetCountThrownUsedProductsInPeriodByUsername(ctx context.Context,
	username string, fromDate string, toDate string) (int, error) {
	var cookedCount int
	err := s.db.GetContext(ctx, &cookedCount, `
	SELECT COUNT(*)
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	WHERE ut.username = $1
	AND f.status = 'thrown'
		AND f.use_date >= $2
		AND f.use_date <= $3;
	`, username, fromDate, toDate)
	if err != nil {
		return -1, err
	}
	return cookedCount, nil
}

func (s *Store) GetTodayBuyList(ctx context.Context) ([]Product, error) {
	products := []Product{}
	err := s.db.SelectContext(ctx, &products, `
	SELECT user_id, id, weight, buy_time, name FROM buy_list
	JOIN product ON product.id = buy_list.product_id
	WHERE buy_time::DATE = CURRENT_DATE ;
	`)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *Store) GetChatIdByUserId(ctx context.Context, userid string) (int64, error) {
	var chatid int64
	err := s.db.GetContext(ctx, &chatid, `
	SELECT chat_id
	FROM usertg
	WHERE id = $1;
	`, userid)
	if err != nil {
		return chatid, err
	}
	return chatid, nil
}
