package bot

import (
	"os"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	ChatID = 1019642784
)

func TestBotCreate(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	bot.Debug = true

	if err != nil {
		t.Error(err)
	}
}

func TestGetUpdate(t *testing.T) {
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)

	_, err := bot.GetUpdates(u)

	if err != nil {
		t.Error(err)
	}
}

func TestSendWithMessage(t *testing.T) {
	bot, _ := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	bot.Debug = true

	msg := tgbotapi.NewMessage(ChatID, "Test send message")
	_, err := bot.Send(msg)

	if err != nil {
		t.Error(err)
	}
}
