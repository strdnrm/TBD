package bot

import (
	"buy_list/bot/store"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	// "strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var (
	p           store.Product
	f           store.FridgeProduct
	ur          store.Usertg
	GlobalState int
	logger      *zap.Logger
	ctx         context.Context
)

func StartBot() {
	logger = zap.NewExample()
	defer logger.Sync()

	ctx = context.Background()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		logger.Panic("Invalid token", zap.Error(err))
		// log.Panic(err)
	}

	bot.Debug = true

	// logger.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	s := GetConn()

	GlobalState = StateStart

	p = store.Product{}
	f = store.FridgeProduct{}
	ur = store.Usertg{}

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			//log
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				msg.ReplyToMessageID = update.Message.MessageID
				switch update.Message.Command() {
				case "start":
					msg.ReplyMarkup = startKeyboard
					msg.Text = "Привет! Я бот, который может управлять вашими покупками и мониторить срок годности продуктов."
					err := s.AddUsertg(ctx, &store.Usertg{
						Username: update.Message.From.UserName,
					})
					if err != nil {
						logger.Error("Adding user error", zap.Error(err))
					}
				case "cancel":
					GlobalState = StateStart
				default:
					msg.Text = "Неверная команда :("
				}
				SendMessage(bot, &msg)

			} else {
				//reply
				//msg.ReplyToMessageID = update.Message.MessageID

				switch GlobalState {

				//start menu
				case StateStart:
					switch update.Message.Text {
					case startKeyboard.Keyboard[0][0].Text:
						GlobalState = StateAddBuyList
						msg.ReplyMarkup = buylistKeyboard
						SendMessage(bot, &msg)
					case startKeyboard.Keyboard[1][0].Text:
						GlobalState = StateAddFridge
						msg.ReplyMarkup = fridgeKeyboard
						SendMessage(bot, &msg)
					case startKeyboard.Keyboard[2][0].Text:
						GlobalState = StateUsedProducts
						msg.ReplyMarkup = usedProductsKeyboard
						SendMessage(bot, &msg)
					case "open":
						msg.ReplyMarkup = startKeyboard
						SendMessage(bot, &msg)
					case "close":
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						SendMessage(bot, &msg)
					}

				//adding to buy list
				case StateAddBuyList:
					switch update.Message.Text {

					case buylistKeyboard.Keyboard[0][0].Text: //add product
						AddProduct(&msg, &update, s)

					case buylistKeyboard.Keyboard[1][0].Text: //get buy list
						ProductList(&msg, &update, s, bot)
						continue

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						GlobalState = StateStart
						msg.ReplyMarkup = startKeyboard

					default:
						AddingToBuyList(&msg, &update, s)

					}
					SendMessage(bot, &msg)

				case StateAddFridge:
					switch update.Message.Text {

					case fridgeKeyboard.Keyboard[0][0].Text: //add product
						AddFridge(&msg, &update, s)

					case fridgeKeyboard.Keyboard[1][0].Text: //get fridge list by alpha
						list, err := s.GetFridgeListByUsernameAlpha(ctx, update.Message.From.UserName)
						if err != nil {
							logger.Error("Getting fridge list by alphabet error", zap.Error(err))
						}
						FridgeList(&msg, &update, s, bot, list)
						continue

					case fridgeKeyboard.Keyboard[2][0].Text: //get fridge list by exp date
						list, err := s.GetFridgeListByUsernameExpDate(ctx, update.Message.From.UserName)
						if err != nil {
							logger.Error("Getting fridge list by exp date error", zap.Error(err))
						}
						FridgeList(&msg, &update, s, bot, list)
						continue

					case fridgeKeyboard.Keyboard[3][0].Text: //cancel
						GlobalState = StateStart
						msg.ReplyMarkup = startKeyboard

					default:
						AddingToFridge(&msg, &update, s)

					}
					SendMessage(bot, &msg)

				case StateUsedProducts:

					switch update.Message.Text {
					case usedProductsKeyboard.Keyboard[0][0].Text:
						GetUsedProdcutsList(&update, s, bot, &msg)

					case usedProductsKeyboard.Keyboard[2][0].Text: //cancel
						GlobalState = StateStart
						msg.ReplyMarkup = startKeyboard
						SendMessage(bot, &msg)
					}

				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			switch update.CallbackQuery.Data {

			case "deleteProductFromBuyList":
				DeleteProductFromBuyList(&update, s, bot)

			case "addToFridgeFromBuyList":
				AddToFridgeFromBuyList(&update, s, bot)

			case "deleteProductFromFridge":
				DeleteProductFromFridge(&update, s, bot)

			case "openProductFromFridge":
				OpenProductFromFridge(&update, s, bot)

			case "setProductCooked":
				SetProductCookedFromFridge(&update, s, bot)

			case "setProductThrown":
				SetProductThrownFromFridge(&update, s, bot)
			}

			if _, err := bot.Request(callback); err != nil { //
				panic(err)
				//log.Fatal(err)
			}
		}
	}
}

func AddProduct(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, s *store.Store) {

	p.State = 0

	msg.Text = "Введите название продукта"
}

func ProductList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	products, err := s.GetBuyListByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Get buy list error", zap.Error(err))
	}
	if len(products) != 0 {
		for i, pr := range products {
			msg.Text = fmt.Sprintf("%d: %s %.2f %s\n", i+1, pr.Name, pr.Weight, pr.BuyDate)
			msg.ReplyMarkup = inlineBuylistKeyboard
			SendMessage(bot, msg)
		}

	} else {
		msg.Text = "Список покупок пуст"
		SendMessage(bot, msg)
	}
}

func AddingToBuyList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, s *store.Store) {
	switch p.State {
	case 0:
		var err error
		p, err = s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		ur, err = s.GetUserByUsername(ctx, update.Message.From.UserName)
		if err != nil {
			logger.Error("Adding prdocut error", zap.Error(err))
		}
		p.UserId = ur.UserId
		msg.Text = "Введите вес/количество"
		p.State = 1

	case 1:
		if _, err := strconv.ParseFloat(update.Message.Text, 64); err != nil {
			msg.Text = "Неверный формат"
		} else {
			p.Weight, _ = strconv.ParseFloat(update.Message.Text, 64)
			msg.Text = "Введите время покупки (YYYY-MM-DD HH:MM)"
			p.State = 2
		}

	case 2:
		ts := update.Message.Text + ":00"
		if _, err := time.Parse("2006-01-02 15:04:05.999", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			p.BuyDate = ts
			err := s.AddProductToBuyList(ctx, &p)
			if err != nil {
				logger.Error("Adding prodcut to buy list error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в список покупок"
			msg.ReplyMarkup = buylistKeyboard
		}
	}
}

func FridgeList(msg *tgbotapi.MessageConfig, update *tgbotapi.Update,
	s *store.Store, bot *tgbotapi.BotAPI, fridgeProducts []store.FridgeProduct) {
	// fridgeProducts := s.GetFridgeListByUsernameAlpha(update.Message.From.UserName)
	if len(fridgeProducts) != 0 {
		for i, pr := range fridgeProducts {
			resText := fmt.Sprintf("%d: %s \n", i+1, pr.Name)
			if pr.Opened {
				resText += "Открыт "
			} else {
				resText += "Не вскрыт "
			}
			// NOW add expire and other dates
			expdate, err := time.Parse("2006-01-02", pr.Expire_date)
			if err != nil {
				panic(err)
				// log.Fatal(err)
			}

			if time.Now().After(expdate) {
				resText += fmt.Sprintf("Просрочен %s ", pr.Expire_date)
			} else {
				resText += fmt.Sprintf("Годен до %s ", pr.Expire_date)
			}

			msg.Text = resText
			msg.ReplyMarkup = inlineFridgeKeyboard
			if _, err := bot.Send(msg); err != nil {
				panic(err)
				// log.Fatal(err)
			}
		}
	} else {
		msg.Text = "Холодильник пуст"
		SendMessage(bot, msg)
	}
}

func AddFridge(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, s *store.Store) {
	var err error
	ur, err = s.GetUserByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Getting user id error", zap.Error(err))
	}
	f.State = 0
	f.UserId = ur.UserId
	msg.Text = "Введите название продукта"
}

func AddingToFridge(msg *tgbotapi.MessageConfig, update *tgbotapi.Update, s *store.Store) {
	switch f.State {
	case 0:
		var err error
		p, err = s.CreateProductByName(ctx, update.Message.Text)
		if err != nil {
			logger.Error("Creating product error", zap.Error(err))
		}
		f.ProductId = p.ProductId
		f.Name = p.Name

		msg.Text = "Укажите срок годности (YYYY-MM-DD)"
		f.State = 1

	case 1:
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			s.AddProductToFridge(ctx, &f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в холодильник"
			msg.ReplyMarkup = fridgeKeyboard
		}

	case 2: //for adding from buy list
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			err := s.AddProductToFridge(ctx, &f)
			if err != nil {
				logger.Error("Adding product to fridge error", zap.Error(err))
			}
			p, err = s.GetProductByName(ctx, f.Name)
			if err != nil {
				logger.Error("Getting product by name error", zap.Error(err))
			}
			err = s.DeleteProductFromBuyListById(ctx, p.ProductId)
			if err != nil {
				logger.Error("Deleting product from but list error", zap.Error(err))
			}
			msg.Text = "Товар добавлен в холодильник"
			msg.ReplyMarkup = buylistKeyboard
			GlobalState = StateAddBuyList
		}

		//TODO
	case 3: //for opening product
		ts := update.Message.Text
		if _, err := time.Parse("2006-01-02", ts); err != nil {
			msg.Text = "Неверный формат"
		} else {
			f.Expire_date = ts
			err = s.OpenProductFromFridgeById(ctx, f.ProductId, f.Expire_date)
			if err != nil {
				logger.Error("Opennig product from fridge error", zap.Error(err))
			}
			msg.Text = "Срок годности изменен"
			msg.ReplyMarkup = fridgeKeyboard
		}
	}
}

func DeleteProductFromBuyList(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	err = s.DeleteProductFromBuyListById(ctx, p.ProductId)
	if err != nil {
		logger.Error("Deleting product from buy list error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из списка покупок", pname))
	SendMessage(bot, &msg)
}

func DeleteProductFromFridge(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	err = s.DeleteProductFromFridgeById(ctx, p.ProductId)
	if err != nil {
		logger.Error("Deleting product error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' удален из холодильника", pname))
	SendMessage(bot, &msg)
}

func OpenProductFromFridge(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	f.ProductId = p.ProductId
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' открыт\nВведите новый срок годности:", pname))
	f.State = 3
	SendMessage(bot, &msg)
}

func SetProductCookedFromFridge(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	p, err := s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	useDate := time.Now().Format("2006-01-02")
	err = s.SetCookedProductFromFridgeById(ctx, p.ProductId, useDate)
	if err != nil {
		logger.Error("Setting cooked prodcut error error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' приготовлен", pname))
	SendMessage(bot, &msg)
}

func SetProductThrownFromFridge(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	p, err := s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Getting product error", zap.Error(err))
	}
	useDate := time.Now().Format("2006-01-02")
	err = s.SetThrownProductFromFridgeById(ctx, p.ProductId, useDate)
	if err != nil {
		logger.Error("Set thorwn porduct error error", zap.Error(err))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		fmt.Sprintf("Продукт '%s' выкинут", pname))
	SendMessage(bot, &msg)
}

func AddToFridgeFromBuyList(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI) {
	pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
	var err error
	p, err = s.GetProductByName(ctx, pname)
	if err != nil {
		logger.Error("Get product error", zap.Error(err))
	}
	ur, err := s.GetUserByUsername(ctx, update.CallbackQuery.From.UserName)
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
	f.State = 2
}

func GetUsedProdcutsList(update *tgbotapi.Update, s *store.Store, bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	list, err := s.GetUsedProductsByUsername(ctx, update.Message.From.UserName)
	if err != nil {
		logger.Error("Get used products list error", zap.Error(err))
	}
	if len(list) != 0 {
		res := ""
		for i, pr := range list {
			res += fmt.Sprintf("%d: %s %s %s \n", i, pr.Name, pr.Status, pr.Use_date)
		}
		msg.Text = res
	} else {
		msg.Text = "Список покупок пуст"

	}
	SendMessage(bot, msg)
}

func SendMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		panic(err)
		//log.Fatal(err)
	}
}

func GetConn() *store.Store {
	s := fmt.Sprintf("postgresql://%s:%s@localhost:5433/?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password"))
	return store.NewStore(s)
}
