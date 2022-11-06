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

//global states of "fsm"
const (
	StateStart int = iota
	StateAddBuyList
	StateAddFridge
	StateOpenProduct
	StateFinalStatus
	StateProductList
	StateUsedProducts
	StateStat
)

type TgBot struct {
	Tkn string `json:"token"`
}

const (
	StateProduct int = iota
	StateWeight
	StateBuyDate
)

const (
	StateFridgeProduct int = iota
	StateFridgeDate
)

type Fridge struct {
	State int
}
