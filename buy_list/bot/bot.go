package bot

import (
	"buy_list/bot/models"
	"buy_list/bot/store"
	"context"
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//go:generate moq -out store_moq_test.go . Storer
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

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error with loading .env file", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
	if err != nil {
		logger.Panic("Invalid token", zap.Error(err))
	}

	store, err := store.NewStore(
		fmt.Sprintf("postgresql://%s:%s@localhost:5433/tgbot?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password")),
	// fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME")),
	)

	if err != nil {
		panic(err)
	}

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

	bot.BotAPI.Debug = true

	logger.Info("Authorized on account", zap.String("name", bot.BotAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.BotAPI.GetUpdatesChan(u)

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
				HandleCommands(&update, &bot, &msg)
			} else {

				switch GlobalState {

				//start menu
				case StateStart:
					StartMenu(&update, &bot, &msg)

				//adding to buy list
				case StateAddBuyList:
					switch update.Message.Text {

					case buylistKeyboard.Keyboard[0][0].Text: //add product
						AddProduct(&update, &msg)

					case buylistKeyboard.Keyboard[1][0].Text: //get buy list
						ProductList(&update, &bot, &msg)
						continue

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToBuyList(&update, &bot, &msg)

					}
					SendMessage(&bot, &msg)

				case StateAddFridge:
					switch update.Message.Text {

					case fridgeKeyboard.Keyboard[0][0].Text: //add product
						AddFridge(&update, &bot, &msg)

					case fridgeKeyboard.Keyboard[1][0].Text: //get fridge list by alpha
						GetFridgeListByUsernameAlphaMenu(&update, &bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[2][0].Text: //get fridge list by exp date
						GetFridgeListByUsernameExpDateMenu(&update, &bot, &msg)
						continue

					case fridgeKeyboard.Keyboard[3][0].Text: //cancel
						CancelMenu(&msg)

					default:
						AddingToFridge(&update, &bot, &msg)

					}
					SendMessage(&bot, &msg)

				case StateUsedProducts:

					HandleStateUserProducts(&update, &bot, &msg)

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
