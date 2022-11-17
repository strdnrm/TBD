package models

type Product struct { // 0 - name ; 1 - weight ; 2 - buydate
	UserId    string `db:"user_id"`
	ProductId string `db:"id"`
	State     int
	Name      string  `db:"name"`
	Weight    float64 `db:"weight"`
	BuyDate   string  `db:"buy_time"`
}

type FridgeProduct struct { // 0 - name ; 1 - expire date ; 2 - open ; 3 - from buy list
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
