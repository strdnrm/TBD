package bot

import (
	"buy_list/bot/store"
	"context"
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var (
	p           store.Product
	f           store.FridgeProduct
	ur          store.Usertg
	ps          store.PeriodStat
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
	}

	bot.Debug = true

	logger.Info("Authorized on account", zap.String("name", bot.Self.UserName))
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	s := store.NewStore(fmt.Sprintf("postgresql://%s:%s@localhost:5433/?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password")))

	GlobalState = StateStart

	p = store.Product{}
	f = store.FridgeProduct{}
	ur = store.Usertg{}
	ps = store.PeriodStat{}

	InitScheduler(s, bot)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if update.Message.IsCommand() {
				msg.ReplyToMessageID = update.Message.MessageID
				switch update.Message.Command() {
				case "start":
					StartUser(&msg, &update, s)
				case "cancel":
					GlobalState = StateStart
				default:
					msg.Text = "Неверная команда :("
				}
				SendMessage(bot, &msg)

			} else {

				switch GlobalState {

				//start menu
				case StateStart:
					StartMenu(&msg, &update, bot)

				//adding to buy list
				case StateAddBuyList:
					switch update.Message.Text {

					case buylistKeyboard.Keyboard[0][0].Text: //add product
						AddProduct(&msg, &update, s)

					case buylistKeyboard.Keyboard[1][0].Text: //get buy list
						ProductList(&msg, &update, s, bot)
						continue

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToBuyList(&msg, &update, s, bot)

					}
					SendMessage(bot, &msg)

				case StateAddFridge:
					switch update.Message.Text {

					case fridgeKeyboard.Keyboard[0][0].Text: //add product
						AddFridge(&msg, &update, s)

					case fridgeKeyboard.Keyboard[1][0].Text: //get fridge list by alpha
						GetFridgeListByUsernameAlphaMenu(&update, s, bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[2][0].Text: //get fridge list by exp date
						GetFridgeListByUsernameExpDateMenu(&update, s, bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[3][0].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToFridge(&msg, &update, s, bot)

					}
					SendMessage(bot, &msg)

				case StateUsedProducts:

					switch update.Message.Text {

					case usedProductsKeyboard.Keyboard[0][0].Text: // get list of use products
						GetAllUsedProducts(&update, s, bot, &msg)
						continue

					case usedProductsKeyboard.Keyboard[1][0].Text:
						msg.Text = "Введите начальную дату (YYYY-MM-DD)"
						ps.State = StateFromDate

					case usedProductsKeyboard.Keyboard[2][0].Text: //cancel
						CancelMenu(&msg)

					default:
						UsedProductStat(&msg, &update, s, bot)

					}
					SendMessage(bot, &msg)

				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			HandleCallbacks(&update, s, bot)

			if _, err := bot.Request(callback); err != nil {
				logger.Panic(err.Error())
			}
		}
	}
}
