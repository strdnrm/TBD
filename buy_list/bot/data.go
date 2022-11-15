package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var startKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список покупок"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Холодильник"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Использованные продукты"),
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

var inlineBuylistKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "deleteProductFromBuyList"),
		tgbotapi.NewInlineKeyboardButtonData("В холодильник", "addToFridgeFromBuyList"),
	),
)

var fridgeKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить продукт в холодильник"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список по алфавиту"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список по сроку годности"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	),
)

var inlineFridgeKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "deleteProductFromFridge"),
		tgbotapi.NewInlineKeyboardButtonData("Открыть", "openProductFromFridge"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Приготовлен", "setProductCooked"),
		tgbotapi.NewInlineKeyboardButtonData("Выкинут", "setProductThrown"),
	),
)

var usedProductsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Продукты"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Статистика за период"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Отмена"),
	),
)

//global states of "fsm"
const (
	StateStart int = iota
	StateAddBuyList
	StateAddFridge
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
	StateFromBuyList
	StateOpen
)

const (
	StateFromDate int = iota
	StateToDate
)
