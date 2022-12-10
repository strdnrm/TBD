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
	"go.uber.org/zap/zapcore"
)

//go:generate moq -out store_moq_test.go . Storer
type Storer interface {
	AddUsertg(ctx context.Context, u *models.Usertg) error
	GetUserByUsername(ctx context.Context, username string) (models.Usertg, error)
	CreateProductByName(ctx context.Context, productName string) (models.Product, error)
	GetProductByName(ctx context.Context, productName string) (models.Product, error)
	DeleteProductFromBuyListById(ctx context.Context, productId string, userid string) error
	DeleteProductFromFridgeById(ctx context.Context, productId string, userid string) error
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
	GetUsedProductsInPeriodByUsername(ctx context.Context, username string, period models.PeriodStat) ([]models.FridgeProduct, error)
	GetCountCookedUsedProductsInPeriodByUsername(ctx context.Context, username string, period models.PeriodStat) (int, error)
	GetCountThrownUsedProductsInPeriodByUsername(ctx context.Context, username string, period models.PeriodStat) (int, error)
	GetTodayBuyList(ctx context.Context) ([]models.Product, error)
	GetChatIdByUserId(ctx context.Context, userid string) (int64, error)
	GetSoonExpireList(ctx context.Context) ([]models.FridgeProduct, error)
}
type Bot struct {
	BotAPI *tgbotapi.BotAPI
	s      Storer
	p      models.Product
	f      models.FridgeProduct
	ur     models.Usertg
	ps     models.PeriodStat
}

var (
	GlobalState int
	logger      *zap.Logger
)

func newBot() *Bot {
	logger = initializeLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error with loading .env file", zap.Error(err))
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TGTOKEN"))
	if err != nil {
		logger.Panic("Invalid token", zap.Error(err))
	}

	store, err := store.NewStore(
		// fmt.Sprintf("postgresql://%s:%s@localhost:5433/tgbot?sslmode=disable", os.Getenv("dbuser"), os.Getenv("password")),
		fmt.Sprintf("postgresql://%s:%s@%s:%s/?sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT")), //, os.Getenv("DBNAME")
	)

	if err != nil {
		panic(err)
	}

	return &Bot{
		BotAPI: bot,
		s:      store,
		p:      models.Product{},
		f:      models.FridgeProduct{},
		ur:     models.Usertg{},
		ps:     models.PeriodStat{},
	}
}

func StartBot() {
	logger = zap.NewExample()
	defer logger.Sync()

	ctx := context.Background()

	bot := newBot()

	bot.BotAPI.Debug = true

	logger.Info("Authorized on account", zap.String("name", bot.BotAPI.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.BotAPI.GetUpdatesChan(u)

	GlobalState = StateStart

	InitScheduler(bot)

	for update := range updates {

		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			if update.Message.IsCommand() {
				bot.HandleCommands(ctx, &update, &msg)
			} else {

				switch GlobalState {
				//start menu
				case StateStart:
					bot.StartMenu(&update, &msg)

				//adding to buy list
				case StateAddBuyList:
					bot.HandleStateBuylist(ctx, &update, &msg)

				//adding to fridge
				case StateAddFridge:
					bot.HandleStateFridge(ctx, &update, &msg)

				//watch stats
				case StateUsedProducts:
					bot.HandleStateUserProducts(ctx, &update, &msg)

				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			bot.HandleCallbacks(ctx, &update)

			if _, err := bot.BotAPI.Request(callback); err != nil {
				logger.Panic(err.Error())
			}
		}
	}
}

func initializeLogger() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logFile, _ := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}
