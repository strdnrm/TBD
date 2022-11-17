package bot

import (
	"buy_list/bot/models"
	"context"
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var token string = os.Getenv("tgtoken")

func TestInitScheduler(t *testing.T) {
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			GetTodayBuyListFunc: func(ctx context.Context) ([]models.Product, error) {
				return []models.Product{}, nil
			},
		},
	}
	InitScheduler(&bot)
	if !schb.IsRunning() {
		t.Error("Buy list alerts not running")
	}
	if !schf.IsRunning() {
		t.Error("Fridge alerts not running")
	}

}

func TestCreateBuyAlerts(t *testing.T) {

}

func TestCreateExpireAlerts(t *testing.T) {

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
