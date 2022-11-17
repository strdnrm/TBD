package bot

import (
	"buy_list/bot/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func StartMenu(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	switch update.Message.Text {
	case startKeyboard.Keyboard[0][0].Text:
		GlobalState = StateAddBuyList
		msg.ReplyMarkup = buylistKeyboard
		SendMessage(bot, msg)
	case startKeyboard.Keyboard[1][0].Text:
		GlobalState = StateAddFridge
		msg.ReplyMarkup = fridgeKeyboard
		SendMessage(bot, msg)
	case startKeyboard.Keyboard[2][0].Text:
		GlobalState = StateUsedProducts
		msg.ReplyMarkup = usedProductsKeyboard
		SendMessage(bot, msg)
	case "open":
		msg.ReplyMarkup = startKeyboard
		SendMessage(bot, msg)
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		SendMessage(bot, msg)
	}
}

func CancelMenu(msg *tgbotapi.MessageConfig) {
	GlobalState = StateStart
	msg.ReplyMarkup = startKeyboard
	fmt.Println(*msg)
}

func StartUser(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	msg.ReplyMarkup = startKeyboard
	msg.Text = "Привет! Я бот, который может управлять вашими покупками и мониторить срок годности продуктов."
	err := bot.s.AddUsertg(ctx, &models.Usertg{
		Username: update.Message.From.UserName,
		ChatId:   update.Message.From.ID,
	})
	if err != nil {
		logger.Error("Adding user error", zap.Error(err))
	}
}

func AddProduct(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) {
	p.State = StateProduct
	msg.Text = "Введите название продукта"
}

func ProductList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	products, err := bot.s.GetBuyListByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Get buy list error", zap.Error(err))
	}
	if len(products) != 0 {
		for i, pr := range products {
			buydate, err := time.Parse(time.RFC3339, pr.BuyDate)
			if err != nil {
				panic(err)
			}
			msg.Text = fmt.Sprintf("%d: %s %.2f %s\n", i+1, pr.Name, pr.Weight, buydate.Format("2006-01-02 15:04"))
			msg.ReplyMarkup = inlineBuylistKeyboard
			SendMessage(bot, msg)
		}

	} else {
		msg.Text = "Список покупок пуст"
		SendMessage(bot, msg)
	}
}

func HandleCallbacks(update *tgbotapi.Update, bot *Bot) {
	switch update.CallbackQuery.Data {

	case "deleteProductFromBuyList":
		DeleteProductFromBuyList(update, bot)

	case "addToFridgeFromBuyList":
		AddToFridgeFromBuyList(update, bot)

	case "deleteProductFromFridge":
		DeleteProductFromFridge(update, bot)

	case "openProductFromFridge":
		OpenProductFromFridge(update, bot)

	case "setProductCooked":
		SetProductCookedFromFridge(update, bot)

	case "setProductThrown":
		SetProductThrownFromFridge(update, bot)
	}
}

func AddingToBuyList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	switch p.State {
	case StateProduct:
		var err error
		p, err = bot.s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		ur, err = bot.s.GetUserByUsername(ctx, update.Message.From.UserName)
		if err != nil {
			logger.Error("Adding prdocut error", zap.Error(err))
		}
		p.UserId = ur.UserId
		msg.Text = "Введите вес/количество"
		p.State = StateWeight

	case StateWeight:
		if _, err := strconv.ParseFloat(update.Message.Text, 64); err != nil {
			msg.Text = "Неверный формат"
		} else {
			p.Weight, _ = strconv.ParseFloat(update.Message.Text, 64)
			msg.Text = "Введите время покупки (YYYY-MM-DD HH:MM)"
			p.State = StateBuyDate
		}

	case StateBuyDate:
		ts := update.Message.Text + ":00"
		if _, err := time.Parse("2006-01-02 15:04:05.999", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			p.BuyDate = ts
			err := bot.s.AddProductToBuyList(ctx, &p)
			if err != nil {
				logger.Error("Adding prodcut to buy list error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в список покупок"
			msg.ReplyMarkup = buylistKeyboard
			UpdateBuyListSchedule(bot)
		}
	}
}

func FridgeList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update,
	bot *Bot, fridgeProducts []models.FridgeProduct) {
	if len(fridgeProducts) != 0 {
		for i, pr := range fridgeProducts {
			resText := fmt.Sprintf("%d: %s \n", i+1, pr.Name)
			if pr.Opened {
				resText += "Открыт "
			} else {
				resText += "Не вскрыт "
			}
			expdate, err := time.Parse(time.RFC3339, pr.Expire_date)

			if err != nil {
				panic(err)
			}

			if time.Now().After(expdate) {
				resText += fmt.Sprintf("Просрочен %s ", expdate.Format("2006-01-02"))
			} else {
				resText += fmt.Sprintf("Годен до %s ", expdate.Format("2006-01-02"))
			}

			msg.Text = resText
			msg.ReplyMarkup = inlineFridgeKeyboard
			SendMessage(bot, msg)
		}
	} else {
		msg.Text = "Холодильник пуст"
		SendMessage(bot, msg)
	}
}

func AddFridge(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	var err error
	ur, err = bot.s.GetUserByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting user id error", zap.Error(err))
	}
	f.State = StateFridgeProduct
	f.UserId = ur.UserId
	msg.Text = "Введите название продукта"
}

func AddingToFridge(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	switch f.State {
	case StateFridgeProduct:
		var err error
		p, err = bot.s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		f.ProductId = p.ProductId
		f.Name = p.Name

		msg.Text = "Укажите срок годности (YYYY-MM-DD)"
		f.State = StateFridgeDate

	case StateFridgeDate:
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			bot.s.AddProductToFridge(ctx, &f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в холодильник"
			msg.ReplyMarkup = fridgeKeyboard
			UpdateExpireSchedule(bot)
		}

	case StateFromBuyList: //for adding from buy list
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			err := bot.s.AddProductToFridge(ctx, &f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			p, err = bot.s.GetProductByName(ctx, f.Name)
			if err != nil {
				logger.Error("Getting product by name error", zap.Error(err))
			}
			err = bot.s.DeleteProductFromBuyListById(ctx, p.ProductId)
			if err != nil {
				logger.Error("Deleting product from but list error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в холодильник"
			msg.ReplyMarkup = buylistKeyboard
			GlobalState = StateAddBuyList
			UpdateExpireSchedule(bot)
			UpdateBuyListSchedule(bot)
		}

		//TODO
	case StateOpen: //for opening product
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			err = bot.s.OpenProductFromFridgeById(ctx, f.ProductId, f.Expire_date)
			if err != nil {
				logger.Error("Opennig product from fridge error", zap.Error(err))
			}
			msg.Text = "Срок годности изменен"
			UpdateExpireSchedule(bot)
			msg.ReplyMarkup = fridgeKeyboard
		}
	}
}

func UsedProductStat(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, bot *Bot) {
	switch ps.State {

	case StateFromDate:
		if ts, err := time.Parse("2006-01-02", update.Message.Text); err != nil {
			msg.Text = "Неверный формат"
		} else {
			ps.FromDate = ts.Format("2006-01-02")
			ps.State = StateToDate
			msg.Text = "Введите конечную дату (YYYY-MM-DD)"
		}

	case StateToDate:
		if ts, err := time.Parse("2006-01-02", update.Message.Text); err != nil {
			msg.Text = "Неверный формат"
		} else {
			fd, err := time.Parse("2006-01-02", ps.FromDate)
			if err != nil {
				panic(err)
			}
			if ts.Before(fd) {
				msg.Text = "Некорректный период. Введите конечную дату (YYYY-MM-DD)"
			} else {
				ps.ToDate = ts.Format("2006-01-02")

				GetPeriodUsedProducts(update, bot, msg, ps)
				cc, err := bot.s.GetCountThrownUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, ps)
				if err != nil {
					logger.Error("Get count used products in period list error", zap.Error(err))
				}

				ct, err := bot.s.GetCountThrownUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, ps)
				if err != nil {
					logger.Error("Get count used products in period list error", zap.Error(err))
				}

				msg.Text = fmt.Sprintf("Выкинуто продуктов: %d\nПриготовлено: %d", cc, ct)
				msg.ReplyMarkup = usedProductsKeyboard
			}

		}
	}
	SendMessage(bot, msg)
}

func DeleteProductFromBuyList(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	err = bot.s.DeleteProductFromBuyListById(ctx, p.ProductId)
	if err != nil {
		logger.Error("Deleting product from buy list error", zap.Error(err))
	}
	UpdateBuyListSchedule(bot)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из списка покупок", pname))
	SendMessage(bot, &msg)
}

func DeleteProductFromFridge(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	err = bot.s.DeleteProductFromFridgeById(ctx, p.ProductId)
	if err != nil {
		logger.Error("Deleting product error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из холодильника", pname))

	SendMessage(bot, &msg)
	UpdateExpireSchedule(bot)
}

func OpenProductFromFridge(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	f.ProductId = p.ProductId
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' открыт\nВведите новый срок годности:", pname))
	f.State = StateOpen

	SendMessage(bot, &msg)
	UpdateExpireSchedule(bot)
}

func SetProductCookedFromFridge(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	p, err := bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	useDate := time.Now().Format("2006-01-02")
	err = bot.s.SetCookedProductFromFridgeById(ctx, p.ProductId, useDate)
	if err != nil {
		logger.Error("Setting cooked prodcut error error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' приготовлен", pname))
	SendMessage(bot, &msg)
	UpdateExpireSchedule(bot)
}

func SetProductThrownFromFridge(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	p, err := bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	useDate := time.Now().Format("2006-01-02")
	err = bot.s.SetThrownProductFromFridgeById(ctx, p.ProductId, useDate)
	if err != nil {
		logger.Error("Set thorwn porduct error error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' выкинут", pname))
	SendMessage(bot, &msg)
	UpdateExpireSchedule(bot)
}

func AddToFridgeFromBuyList(update *tgbotapi.Update, bot *Bot) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Get product error", zap.Error(err))
	}
	ur, err := bot.s.GetUserByUsername(ctx, update.CallbackQuery.From.UserName)
	if err != nil {
		logger.Error("Get user error", zap.Error(err))
	}
	f.Name = pname
	f.UserId = ur.UserId
	f.ProductId = p.ProductId
	f.Opened = false
	GlobalState = StateAddFridge
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Укажите срок годности (YYYY-MM-DD)")
	SendMessage(bot, &msg)
	f.State = StateFromBuyList
}

func GetAllUsedProducts(update *tgbotapi.Update, bot *Bot, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetUsedProductsByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Get used products list error", zap.Error(err))
	}
	GetUsedProdcutsList(update, bot, msg, list)
}

func GetPeriodUsedProducts(update *tgbotapi.Update, bot *Bot,
	msg *tgbotapi.MessageConfig, period models.PeriodStat) {
	list, err := bot.s.GetUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, period)
	if err != nil {
		logger.Error("Get used products in period list error", zap.Error(err))
	}
	GetUsedProdcutsList(update, bot, msg, list)
}

func GetUsedProdcutsList(update *tgbotapi.Update, bot *Bot, msg *tgbotapi.MessageConfig, list []models.FridgeProduct) {
	if len(list) != 0 {
		res := ""
		for i, pr := range list {
			useDate, err := time.Parse(time.RFC3339, pr.Use_date)
			if err != nil {
				panic(err)
			}
			var st string
			switch pr.Status {
			case "cooked":
				st = "Приготовлен"
			case "thrown":
				st = "Выкинут"
			}
			res += fmt.Sprintf("%d: %s %s %s \n", i+1, pr.Name, st, useDate.Format("2006-01-02"))
		}
		msg.Text = res
	} else {
		msg.Text = "Нет использованных продуктов"

	}
	SendMessage(bot, msg)
}

func GetFridgeListByUsernameAlphaMenu(update *tgbotapi.Update, bot *Bot, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetFridgeListByUsernameAlpha(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting fridge list by alphabet error", zap.Error(err))
	}
	FridgeList(msg, update, bot, list)
}

func GetFridgeListByUsernameExpDateMenu(update *tgbotapi.Update, bot *Bot, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetFridgeListByUsernameExpDate(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting fridge list by exp date error", zap.Error(err))
	}
	FridgeList(msg, update, bot, list)
}

func SendMessage(bot *Bot, msg *tgbotapi.MessageConfig) {
	if _, err := bot.BotAPI.Send(msg); err != nil {
		panic(err)
		//log.Fatal(err)
	}
}
