package bot

import (
	"buy_list/bot/models"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func (bot *Bot) StartMenu(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Text {
	case startKeyboard.Keyboard[0][0].Text:
		GlobalState = StateAddBuyList
		msg.ReplyMarkup = buylistKeyboard
		bot.SendMessage(msg)
	case startKeyboard.Keyboard[1][0].Text:
		GlobalState = StateAddFridge
		msg.ReplyMarkup = fridgeKeyboard
		bot.SendMessage(msg)
	case startKeyboard.Keyboard[2][0].Text:
		GlobalState = StateUsedProducts
		msg.ReplyMarkup = usedProductsKeyboard
		bot.SendMessage(msg)
	case "open":
		msg.ReplyMarkup = startKeyboard
		bot.SendMessage(msg)
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		bot.SendMessage(msg)
	}
}

func (bot *Bot) HandleCommands(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	msg.ReplyToMessageID = update.Message.MessageID
	switch update.Message.Command() {
	case "start":
		bot.StartUser(ctx, update, msg)
	default:
		msg.Text = "Неверная команда :("
	}
	bot.SendMessage(msg)
}

func (bot *Bot) HandleStateBuylist(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Text {

	case buylistKeyboard.Keyboard[0][0].Text: //add product
		bot.AddProduct(update, msg)

	case buylistKeyboard.Keyboard[1][0].Text: //get buy list
		bot.ProductList(ctx, update, msg)

	case buylistKeyboard.Keyboard[1][1].Text: //cancel
		bot.CancelMenu(msg)

	default:
		if bot.p.State != StateProductNull {
			bot.AddingToBuyList(ctx, update, msg)
		} else {
			msg.Text = ""
		}

	}
	bot.SendMessage(msg)
}

func (bot *Bot) HandleStateFridge(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Text {

	case fridgeKeyboard.Keyboard[0][0].Text: //add product
		bot.AddFridge(ctx, update, msg)

	case fridgeKeyboard.Keyboard[1][0].Text: //get fridge list by alpha
		bot.GetFridgeListByUsernameAlphaMenu(ctx, update, msg)

	case fridgeKeyboard.Keyboard[2][0].Text: //get fridge list by exp date
		bot.GetFridgeListByUsernameExpDateMenu(ctx, update, msg)

	case fridgeKeyboard.Keyboard[3][0].Text: //cancel
		bot.CancelMenu(msg)

	default:
		if bot.f.State != StateFridgeNull {
			bot.AddingToFridge(ctx, update, msg)
		} else {
			msg.Text = ""
		}

	}
	bot.SendMessage(msg)
}

func (bot *Bot) HandleStateUserProducts(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch update.Message.Text {

	case usedProductsKeyboard.Keyboard[0][0].Text: // get list of use products
		bot.GetAllUsedProducts(ctx, update, msg)

	case usedProductsKeyboard.Keyboard[1][0].Text:
		bot.GetToDate(msg)

	case usedProductsKeyboard.Keyboard[2][0].Text: //cancel
		bot.CancelMenu(msg)

	default:
		if bot.ps.State != StateDateNull {
			bot.UsedProductStat(ctx, msg, update)
		} else {
			msg.Text = ""
		}

	}
	bot.SendMessage(msg)
}

func (bot *Bot) CancelMenu(msg *tgbotapi.MessageConfig) {
	GlobalState = StateStart
	msg.ReplyMarkup = startKeyboard
	if bot.p.State != StateProductNull {
		bot.p.State = StateProductNull
	}
	if bot.f.State != StateFridgeNull {
		bot.f.State = StateFridgeNull
	}
	if bot.ps.State != StateDateNull {
		bot.ps.State = StateDateNull
	}
}

func (bot *Bot) StartUser(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
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

func (bot *Bot) AddProduct(update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	bot.p.State = StateProduct
	msg.Text = "Введите название продукта"
}

func (bot *Bot) ProductList(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
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
			bot.SendMessage(msg)
		}
		msg.Text = ""

	} else {
		msg.Text = "Список покупок пуст"
	}
}

func (bot *Bot) HandleCallbacks(ctx context.Context, update *tgbotapi.Update) {
	switch update.CallbackQuery.Data {

	case "deleteProductFromBuyList":
		bot.DeleteProductFromBuyList(ctx, update)

	case "addToFridgeFromBuyList":
		bot.AddToFridgeFromBuyList(ctx, update)

	case "deleteProductFromFridge":
		bot.DeleteProductFromFridge(ctx, update)

	case "openProductFromFridge":
		bot.OpenProductFromFridge(ctx, update)

	case "setProductCooked":
		bot.SetProductCookedFromFridge(ctx, update)

	case "setProductThrown":
		bot.SetProductThrownFromFridge(ctx, update)
	}
}

func (bot *Bot) AddingToBuyList(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch bot.p.State {
	case StateProduct:
		var err error
		bot.p, err = bot.s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		bot.ur, err = bot.s.GetUserByUsername(ctx, update.Message.From.UserName)
		if err != nil {
			logger.Error("Adding prdocut error", zap.Error(err))
		}
		bot.p.UserId = bot.ur.UserId
		msg.Text = "Введите вес/количество"
		bot.p.State = StateWeight

	case StateWeight:
		if _, err := strconv.ParseFloat(update.Message.Text, 64); err != nil {
			msg.Text = "Неверный формат"
		} else {
			bot.p.Weight, _ = strconv.ParseFloat(update.Message.Text, 64)
			msg.Text = "Введите время покупки (YYYY-MM-DD HH:MM)"
			bot.p.State = StateBuyDate
		}

	case StateBuyDate:
		ts := update.Message.Text + ":00"
		if _, err := time.Parse("2006-01-02 15:04:05.999", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			bot.p.BuyDate = ts
			err := bot.s.AddProductToBuyList(ctx, &bot.p)
			if err != nil {
				logger.Error("Adding prodcut to buy list error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в список покупок"
			msg.ReplyMarkup = buylistKeyboard
			bot.p.State = StateProductNull
			UpdateBuyListSchedule(bot)
		}
	}
}

func (bot *Bot) FridgeList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, fridgeProducts []models.FridgeProduct) {
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
			bot.SendMessage(msg)
		}
		msg.Text = ""
	} else {
		msg.Text = "Холодильник пуст"
	}
}

func (bot *Bot) AddFridge(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	var err error
	bot.ur, err = bot.s.GetUserByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting user id error", zap.Error(err))
	}
	bot.f.State = StateFridgeProduct
	bot.f.UserId = bot.ur.UserId
	msg.Text = "Введите название продукта"
}

func (bot *Bot) AddingToFridge(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	switch bot.f.State {
	case StateFridgeProduct:
		var err error
		bot.p, err = bot.s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		bot.f.ProductId = bot.p.ProductId
		bot.f.Name = bot.p.Name

		msg.Text = "Укажите срок годности (YYYY-MM-DD)"
		bot.f.State = StateFridgeDate

	case StateFridgeDate:
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			bot.f.Expire_date = ts
			bot.s.AddProductToFridge(ctx, &bot.f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в холодильник"
			msg.ReplyMarkup = fridgeKeyboard
			bot.f.State = StateFridgeNull
			UpdateExpireSchedule(bot)
		}

	case StateFromBuyList: //for adding from buy list
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			bot.f.Expire_date = ts
			err := bot.s.AddProductToFridge(ctx, &bot.f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			bot.p, err = bot.s.GetProductByName(ctx, bot.f.Name)
			if err != nil {
				logger.Error("Getting product by name error", zap.Error(err))
			}
			u, err := bot.s.GetUserByUsername(context.Background(), update.Message.From.UserName)
			if err != nil {
				logger.Error("Getting user error", zap.Error(err))
			}

			err = bot.s.DeleteProductFromBuyListById(ctx, bot.p.ProductId, u.UserId)
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
			bot.f.Expire_date = ts
			err = bot.s.OpenProductFromFridgeById(ctx, bot.f.ProductId, bot.f.Expire_date)
			if err != nil {
				logger.Error("Opennig product from fridge error", zap.Error(err))
			}
			msg.Text = "Срок годности изменен"
			UpdateExpireSchedule(bot)
			msg.ReplyMarkup = fridgeKeyboard
		}
	}
}

func (bot *Bot) GetToDate(msg *tgbotapi.MessageConfig) {
	msg.Text = "Введите начальную дату (YYYY-MM-DD)"
	bot.ps.State = StateFromDate
}

func (bot *Bot) UsedProductStat(ctx context.Context, msg *tgbotapi.MessageConfig, update *tgbotapi.Update) {
	switch bot.ps.State {

	case StateFromDate:
		if ts, err := time.Parse("2006-01-02", update.Message.Text); err != nil {
			msg.Text = "Неверный формат"
		} else {
			bot.ps.FromDate = ts.Format("2006-01-02")
			bot.ps.State = StateToDate
			msg.Text = "Введите конечную дату (YYYY-MM-DD)"
		}

	case StateToDate:
		if ts, err := time.Parse("2006-01-02", update.Message.Text); err != nil {
			msg.Text = "Неверный формат"
		} else {
			fd, err := time.Parse("2006-01-02", bot.ps.FromDate)
			if err != nil {
				panic(err)
			}
			if ts.Before(fd) {
				msg.Text = "Некорректный период. Введите конечную дату (YYYY-MM-DD)"
			} else {
				bot.ps.ToDate = ts.Format("2006-01-02")

				bot.GetPeriodUsedProducts(ctx, update, msg, bot.ps)
				cc, err := bot.s.GetCountThrownUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, bot.ps)
				if err != nil {
					logger.Error("Get count used products in period list error", zap.Error(err))
				}

				ct, err := bot.s.GetCountCookedUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, bot.ps)
				if err != nil {
					logger.Error("Get count used products in period list error", zap.Error(err))
				}
				bot.ps.State = StateDateNull
				msg.Text = fmt.Sprintf("Выкинуто продуктов: %d\nПриготовлено: %d", cc, ct)
				msg.ReplyMarkup = usedProductsKeyboard
			}

		}
	}
}

func (bot *Bot) DeleteProductFromBuyList(ctx context.Context, update *tgbotapi.Update) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	bot.p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	u, err := bot.s.GetUserByUsername(context.Background(), update.CallbackQuery.From.UserName)
	if err != nil {
		logger.Error("Getting user error", zap.Error(err))
	}
	err = bot.s.DeleteProductFromBuyListById(ctx, bot.p.ProductId, u.UserId)
	if err != nil {
		logger.Error("Deleting product from buy list error", zap.Error(err))
	}
	UpdateBuyListSchedule(bot)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из списка покупок", pname))
	bot.SendMessage(&msg)
}

func (bot *Bot) DeleteProductFromFridge(ctx context.Context, update *tgbotapi.Update) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	bot.p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	u, err := bot.s.GetUserByUsername(context.Background(), update.CallbackQuery.From.UserName)
	if err != nil {
		logger.Error("Getting user error", zap.Error(err))
	}
	err = bot.s.DeleteProductFromFridgeById(ctx, bot.p.ProductId, u.UserId)
	if err != nil {
		logger.Error("Deleting product error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из холодильника", pname))

	bot.SendMessage(&msg)
	UpdateExpireSchedule(bot)
}

func (bot *Bot) OpenProductFromFridge(ctx context.Context, update *tgbotapi.Update) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	bot.p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	bot.f.ProductId = bot.p.ProductId
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' открыт\nВведите новый срок годности:", pname))
	bot.f.State = StateOpen

	bot.SendMessage(&msg)
	UpdateExpireSchedule(bot)
}

func (bot *Bot) SetProductCookedFromFridge(ctx context.Context, update *tgbotapi.Update) {
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
	bot.SendMessage(&msg)
	UpdateExpireSchedule(bot)
}

func (bot *Bot) SetProductThrownFromFridge(ctx context.Context, update *tgbotapi.Update) {
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
	bot.SendMessage(&msg)
	UpdateExpireSchedule(bot)
}

func (bot *Bot) AddToFridgeFromBuyList(ctx context.Context, update *tgbotapi.Update) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	bot.p, err = bot.s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Get product error", zap.Error(err))
	}
	ur, err := bot.s.GetUserByUsername(ctx, update.CallbackQuery.From.UserName)
	if err != nil {
		logger.Error("Get user error", zap.Error(err))
	}
	bot.f.Name = pname
	bot.f.UserId = ur.UserId
	bot.f.ProductId = bot.p.ProductId
	bot.f.Opened = false
	GlobalState = StateAddFridge
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Укажите срок годности (YYYY-MM-DD)")
	bot.SendMessage(&msg)
	bot.f.State = StateFromBuyList
}

func (bot *Bot) GetAllUsedProducts(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetUsedProductsByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Get used products list error", zap.Error(err))
	}
	bot.GetUsedProdcutsList(update, msg, list)
}

func (bot *Bot) GetPeriodUsedProducts(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig, period models.PeriodStat) {
	list, err := bot.s.GetUsedProductsInPeriodByUsername(ctx, update.Message.From.UserName, period)
	if err != nil {
		logger.Error("Get used products in period list error", zap.Error(err))
	}
	bot.GetUsedProdcutsList(update, msg, list)
}

func (bot *Bot) GetUsedProdcutsList(update *tgbotapi.Update, msg *tgbotapi.MessageConfig, list []models.FridgeProduct) {
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
}

func (bot *Bot) GetFridgeListByUsernameAlphaMenu(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetFridgeListByUsernameAlpha(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting fridge list by alphabet error", zap.Error(err))
	}
	bot.FridgeList(msg, update, list)
}

func (bot *Bot) GetFridgeListByUsernameExpDateMenu(ctx context.Context, update *tgbotapi.Update, msg *tgbotapi.MessageConfig) {
	list, err := bot.s.GetFridgeListByUsernameExpDate(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting fridge list by exp date error", zap.Error(err))
	}
	bot.FridgeList(msg, update, list)
}

func (bot *Bot) SendMessage(msg *tgbotapi.MessageConfig) {
	if msg.Text != "" {
		if _, err := bot.BotAPI.Send(msg); err != nil {
			panic(err)
		}
	}
}
