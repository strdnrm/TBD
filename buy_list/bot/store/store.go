package store

import (
	"buy_list/bot/models"
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	db *sqlx.DB
}

//TODO inteface

func NewStore(connString string) (*Store, error) {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:./migrations/",
		"postgres", driver)
	if err != nil {
		return nil, err
	}
	m.Up()

	return &Store{
		db: db,
	}, nil
}

func (s *Store) AddUsertg(ctx context.Context, u *models.Usertg) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO usertg(username, chat_id)
	VALUES ($1, $2);
	`, u.Username, u.ChatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (models.Usertg, error) {
	u := models.Usertg{}
	err := s.db.GetContext(ctx, &u, `
	SELECT id::text FROM usertg WHERE username = $1;
	`, username)
	if err != nil {
		return u, err
	}
	return u, nil
}

// returns id if the product exists otherwise creates it
func (s *Store) CreateProductByName(ctx context.Context, productName string) (models.Product, error) {
	p := models.Product{}

	err := s.db.GetContext(ctx, &p, `
	WITH s AS (
		SELECT id, name
		FROM product
		WHERE name = $1
	), i AS (
		INSERT INTO product(name)
		SELECT $1
		WHERE NOT EXISTS (SELECT 1 FROM s)
		RETURNING id
	)
	SELECT id
	FROM i
	UNION ALL
	SELECT id
	FROM s
	`, productName)
	p.Name = productName
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s *Store) GetProductByName(ctx context.Context, productName string) (models.Product, error) {
	p := models.Product{}
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

func (s *Store) AddProductToBuyList(ctx context.Context, p *models.Product) error {
	_, err := s.db.ExecContext(ctx, `
	INSERT INTO buy_list
	VALUES($1, $2, $3, $4);
	`, p.UserId, p.ProductId, p.Weight, p.BuyDate)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBuyListByUsername(ctx context.Context, username string) ([]models.Product, error) {
	// get name wight buydate
	var list []models.Product
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

func (s *Store) AddProductToFridge(ctx context.Context, f *models.FridgeProduct) error {
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

func (s *Store) GetFridgeListByUsername(ctx context.Context, username string) ([]models.FridgeProduct, error) {
	var list []models.FridgeProduct
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

func (s *Store) GetFridgeListByUsernameAlpha(ctx context.Context, username string) ([]models.FridgeProduct, error) {
	var list []models.FridgeProduct
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

func (s *Store) GetFridgeListByUsernameExpDate(ctx context.Context, username string) ([]models.FridgeProduct, error) {
	//get name opened expire_date  status
	var list []models.FridgeProduct
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

func (s *Store) GetUsedProductsByUsername(ctx context.Context, username string) ([]models.FridgeProduct, error) {
	products := []models.FridgeProduct{}
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
	username string, period models.PeriodStat) ([]models.FridgeProduct, error) {
	products := []models.FridgeProduct{}
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
	`, username, period.FromDate, period.ToDate)
	if err != nil {
		return products, err
	}
	return products, nil
}

func (s *Store) GetCountCookedUsedProductsInPeriodByUsername(ctx context.Context,
	username string, period models.PeriodStat) (int, error) {
	var cookedCount int
	err := s.db.GetContext(ctx, &cookedCount, `
	SELECT COUNT(*)
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	WHERE ut.username = $1
	AND f.status = 'cooked'
		AND f.use_date >= $2
		AND f.use_date <= $3;
	`, username, period.FromDate, period.ToDate)
	if err != nil {
		return -1, err
	}
	return cookedCount, nil
}

func (s *Store) GetCountThrownUsedProductsInPeriodByUsername(ctx context.Context,
	username string, period models.PeriodStat) (int, error) {
	var cookedCount int
	err := s.db.GetContext(ctx, &cookedCount, `
	SELECT COUNT(*)
	FROM fridge f
	JOIN usertg ut ON ut.id = f.user_id
	WHERE ut.username = $1
	AND f.status = 'thrown'
		AND f.use_date >= $2
		AND f.use_date <= $3;
	`, username, period.FromDate, period.ToDate)
	if err != nil {
		return -1, err
	}
	return cookedCount, nil
}

func (s *Store) GetTodayBuyList(ctx context.Context) ([]models.Product, error) {
	products := []models.Product{}
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

func (s *Store) GetSoonExpireList(ctx context.Context) ([]models.FridgeProduct, error) {
	products := []models.FridgeProduct{}
	err := s.db.SelectContext(ctx, &products, `
	SELECT user_id, product_id, name, expire_date FROM fridge
	JOIN product ON product.id = fridge.product_id
	WHERE (expire_date = CURRENT_DATE 
	OR expire_date = CURRENT_DATE + INTERVAL '1 day')
	AND status IS NULL
	`)
	if err != nil {
		return products, err
	}
	return products, nil
}
