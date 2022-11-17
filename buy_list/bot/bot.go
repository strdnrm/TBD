package bot

import (
	"buy_list/bot/models"
	"buy_list/bot/store"
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Storer interface {
	AddUsertg(ctx context.Context, u *models.Usertg) error
	GetUserByUsername(ctx context.Context, username string) (models.Usertg, error)
	CreateProductByName(ctx context.Context, productName string) (models.Product, error)
	GetProductByName(ctx context.Context, productName string) (models.Product, error)
	DeleteProductFromBuyListById(ctx context.Context, productId string) error
	DeleteProductFromFridgeById(ctx context.Context, productId string) error
	OpenProductFromFridgeById(ctx context.Context, productId string, expDate string) error
	SetCookedProductFromFridgeById(ctx context.Context, productId string, useDate string) error
	SetThrownProductFromFridgeById(ctx context.Context, productId string, useDate string) error
	AddProductToBuyList(ctx context.Context, p *models.Product) error
	GetBuyListByUsername(ctx context.Context, username string) ([]models.Product, error)
	AddProductToFridge(ctx context.Context, f *models.FridgeProduct) error
	GetFridgeListByUsername(ctx context.Context, username string) ([]models.FridgeProduct, error)
	GetFridgeListByUsernameAlpha(ctx context.Context, username string) ([]models.FridgeProduct, error)
	GetFridgeListByUsernameExpDate(ctx context.Context, username string) ([]models.FridgeProduct, error)
	GetUsedProductsByUsername(ctx context.Context, username string) ([]models.FridgeProduct, error)
	GetUsedProductsInPeriodByUsername(ctx context.Context,
		username string, period models.PeriodStat) ([]models.FridgeProduct, error)
	GetCountCookedUsedProductsInPeriodByUsername(ctx context.Context,
		username string, period models.PeriodStat) (int, error)
	GetCountThrownUsedProductsInPeriodByUsername(ctx context.Context,
		username string, period models.PeriodStat) (int, error)
	GetTodayBuyList(ctx context.Context) ([]models.Product, error)
	GetChatIdByUserId(ctx context.Context, userid string) (int64, error)
	GetSoonExpireList(ctx context.Context) ([]models.FridgeProduct, error)
}

type Bot struct {
	BotAPI *tgbotapi.BotAPI
	s      Storer
}

var (
	p           models.Product
	f           models.FridgeProduct
	ur          models.Usertg
	ps          models.PeriodStat
	GlobalState int
	logger      *zap.Logger
	ctx         context.Context
)

func newBot() Bot {
	logger = zap.NewExample()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		logger.Panic("Invalid token", zap.Error(err))
	}

	store := store.NewStore(fmt.Sprintf("postgresql://%s:%s@localhost:5433/?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password")))

	return Bot{
		BotAPI: bot,
		s:      store,
	}
}

func StartBot() {
	logger = zap.NewExample()
	defer logger.Sync()

	ctx = context.Background()

	bot := newBot()

	// bot, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	// if err != nil {
	// 	logger.Panic("Invalid token", zap.Error(err))
	// }

	bot.BotAPI.Debug = true

	logger.Info("Authorized on account", zap.String("name", bot.BotAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.BotAPI.GetUpdatesChan(u)

	// s := store.NewStore(fmt.Sprintf("postgresql://%s:%s@localhost:5433/?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password")))

	GlobalState = StateStart

	p = models.Product{}
	f = models.FridgeProduct{}
	ur = models.Usertg{}
	ps = models.PeriodStat{}

	InitScheduler(&bot)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			if update.Message.IsCommand() {
				msg.ReplyToMessageID = update.Message.MessageID
				switch update.Message.Command() {
				case "start":
					StartUser(&msg, &update, &bot)
				case "cancel":
					GlobalState = StateStart
				default:
					msg.Text = "Неверная команда :("
				}
				SendMessage(&bot, &msg)

			} else {

				switch GlobalState {

				//start menu
				case StateStart:
					StartMenu(&msg, &update, &bot)

				//adding to buy list
				case StateAddBuyList:
					switch update.Message.Text {

					case buylistKeyboard.Keyboard[0][0].Text: //add product
						AddProduct(&msg, &update)

					case buylistKeyboard.Keyboard[1][0].Text: //get buy list
						ProductList(&msg, &update, &bot)
						continue

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToBuyList(&msg, &update, &bot)

					}
					SendMessage(&bot, &msg)

				case StateAddFridge:
					switch update.Message.Text {

					case fridgeKeyboard.Keyboard[0][0].Text: //add product
						AddFridge(&msg, &update, &bot)

					case fridgeKeyboard.Keyboard[1][0].Text: //get fridge list by alpha
						GetFridgeListByUsernameAlphaMenu(&update, &bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[2][0].Text: //get fridge list by exp date
						GetFridgeListByUsernameExpDateMenu(&update, &bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[3][0].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToFridge(&msg, &update, &bot)

					}
					SendMessage(&bot, &msg)

				case StateUsedProducts:

					switch update.Message.Text {

					case usedProductsKeyboard.Keyboard[0][0].Text: // get list of use products
						GetAllUsedProducts(&update, &bot, &msg)
						continue

					case usedProductsKeyboard.Keyboard[1][0].Text:
						msg.Text = "Введите начальную дату (YYYY-MM-DD)"
						ps.State = StateFromDate

					case usedProductsKeyboard.Keyboard[2][0].Text: //cancel
						CancelMenu(&msg)

					default:
						UsedProductStat(&msg, &update, &bot)

					}
					SendMessage(&bot, &msg)

				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			HandleCallbacks(&update, &bot)

			if _, err := bot.BotAPI.Request(callback); err != nil {
				logger.Panic(err.Error())
			}
		}
	}
}
