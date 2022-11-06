package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список покупок"),
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

var buylistKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить товар в список покупок"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список"),
		tgbotapi.NewKeyboardButton("Отмена"),
	),
)

var deleteKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "datahet"),
		tgbotapi.NewInlineKeyboardButtonData("Добавить в холодильник", "addfirdge"),
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
