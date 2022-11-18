package bot

import (
	"buy_list/bot/models"
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func TestInitScheduler(t *testing.T) {
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			GetTodayBuyListFunc: func(ctx context.Context) ([]models.Product, error) {
				return []models.Product{
					{
						UserId:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:     0,
						Name:      "Pillow",
						Weight:    10.0,
						BuyDate:   "2006-01-02T15:04:05Z",
					},
					{
						UserId:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:     0,
						Name:      "Mango",
						Weight:    10.0,
						BuyDate:   "2006-01-02T15:04:05Z",
					},
				}, nil
			},
			GetChatIdByUserIdFunc: func(ctx context.Context, userid string) (int64, error) {
				return 1019642784, nil
			},
			GetSoonExpireListFunc: func(ctx context.Context) ([]models.FridgeProduct, error) {
				return []models.FridgeProduct{
					{
						UserId:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:       0,
						Name:        "Pillow",
						Opened:      false,
						Expire_date: "2006-01-02T15:04:05Z",
					},
					{
						UserId:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:       0,
						Name:        "Mango",
						Opened:      false,
						Expire_date: "2006-01-02T15:04:05Z",
					},
				}, nil
			},
		},
	}
	logger = zap.NewExample()
	defer logger.Sync()
	InitScheduler(&bot)
	if !schb.IsRunning() {
		t.Error("Buy list alerts not running")
	}
	if !schf.IsRunning() {
		t.Error("Fridge alerts not running")
	}

}

func TestCreateBuyAlerts(t *testing.T) {
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			GetTodayBuyListFunc: func(ctx context.Context) ([]models.Product, error) {
				return []models.Product{
					{
						UserId:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:     0,
						Name:      "Pillow",
						Weight:    10.0,
						BuyDate:   time.Now().Add(time.Hour).Format(time.RFC3339),
					},
					{
						UserId:    "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId: "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:     0,
						Name:      "Mango",
						Weight:    10.0,
						BuyDate:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
					},
				}, nil
			},
			GetChatIdByUserIdFunc: func(ctx context.Context, userid string) (int64, error) {
				return 1019642784, nil
			},
		},
	}
	logger = zap.NewExample()
	defer logger.Sync()
	schb = gocron.NewScheduler(time.Local)
	CreateBuyAlerts(&bot)
	if len(schb.Jobs()) != 2 {
		t.Errorf("Want %d jobs got %d", 2, len(schb.Jobs()))
	}

}

func TestCreateExpireAlerts(t *testing.T) {
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			GetSoonExpireListFunc: func(ctx context.Context) ([]models.FridgeProduct, error) {
				return []models.FridgeProduct{
					{
						UserId:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:       0,
						Name:        "Pillow",
						Opened:      false,
						Expire_date: time.Now().Add(time.Hour).Format(time.RFC3339),
					},
					{
						UserId:      "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						ProductId:   "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
						State:       0,
						Name:        "Mango",
						Opened:      false,
						Expire_date: time.Now().Add(time.Hour).Format(time.RFC3339),
					},
				}, nil
			},
			GetChatIdByUserIdFunc: func(ctx context.Context, userid string) (int64, error) {
				return 1019642784, nil
			},
		},
	}
	logger = zap.NewExample()
	defer logger.Sync()
	schf = gocron.NewScheduler(time.Local)
	CreateExpireAlerts(&bot)
	// zero because we send notifications immediately
	if len(schf.Jobs()) != 0 {
		t.Errorf("Want %d jobs got %d", 0, len(schf.Jobs()))
	}

}

func TestSendAlert(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	msg := tgbotapi.NewMessage(ChatID, "Test send alert")
	if _, err := bot.Send(msg); err != nil {
		t.Error(err)
	}
}
