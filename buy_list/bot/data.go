package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить в список покупок"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить в холодильник"),
		tgbotapi.NewKeyboardButton("Открыть продукт"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Финальный статус продукта"),
		tgbotapi.NewKeyboardButton("Просмотреть список продуктов"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Просмотреть использованных продуктов"),
		tgbotapi.NewKeyboardButton("Статистика"),
	),
)

type TgBot struct {
	Tkn string `json:"token"`
}

const (
	StateProduct int = iota
	StateWeight
	StateBuyDate
)

type Product struct { // 0 - name ; 1 - weight ; 2 - buydate
	State   int
	Name    string
	Weight  float64
	BuyDate string
}
