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

func TestStartMenuBuy(t *testing.T) {
	update := tgbotapi.Update{
		UpdateID: 404723786,
		Message: &tgbotapi.Message{
			Text: startKeyboard.Keyboard[0][0].Text,
		},
	}
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s:      &StorerMock{},
	}
	msg := tgbotapi.NewMessage(ChatID, "test start menu")
	bot.StartMenu(&update, &msg)
	if GlobalState != StateAddBuyList {
		t.Error("Buy list state disabled")
	}

	GlobalState = StateStart
	update = tgbotapi.Update{
		UpdateID: 404723786,
		Message: &tgbotapi.Message{
			Text: startKeyboard.Keyboard[1][0].Text,
		},
	}
	bot.StartMenu(&update, &msg)
	if GlobalState != StateAddFridge {
		t.Error("Fridge state disabled")
	}

	GlobalState = StateStart
	update = tgbotapi.Update{
		UpdateID: 404723786,
		Message: &tgbotapi.Message{
			Text: startKeyboard.Keyboard[2][0].Text,
		},
	}
	bot.StartMenu(&update, &msg)
	if GlobalState != StateUsedProducts {
		t.Error("State used products state disabled")
	}
}

func TestCancelMenu(t *testing.T) {
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s:      &StorerMock{},
		p:      models.Product{},
		f:      models.FridgeProduct{},
		ur:     models.Usertg{},
		ps:     models.PeriodStat{},
	}

	GlobalState = StateAddBuyList
	msg := tgbotapi.NewMessage(ChatID, "test start menu")
	bot.CancelMenu(&msg)
	if GlobalState != StateStart {
		t.Error("State start disabled")
	}
}

func TestStartUser(t *testing.T) {

	update := tgbotapi.Update{
		UpdateID: 404723786,
		Message: &tgbotapi.Message{
			Text: "aaa",
			From: &tgbotapi.User{
				UserName: "name",
				ID:       404723786,
			},
		},
	}
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			AddUsertgFunc: func(ctx context.Context, u *models.Usertg) error {
				return nil
			},
		},
	}

	msg := tgbotapi.NewMessage(ChatID, "test start menu")
	bot.StartUser(context.Background(), &update, &msg)
	// how
}

func TestHandleCallbacks(t *testing.T) {
	update := tgbotapi.Update{
		UpdateID: 404723786,
		CallbackQuery: &tgbotapi.CallbackQuery{
			ID:   "411111111",
			Data: "deleteProductFromBuyList",
			Message: &tgbotapi.Message{
				Text: "1: aaa",
				Chat: &tgbotapi.Chat{ID: ChatID},
			},
		},
	}
	b, err := tgbotapi.NewBotAPI(os.Getenv("tgtoken"))
	if err != nil {
		t.Error(err)
	}
	bot := Bot{
		BotAPI: b,
		s: &StorerMock{
			GetProductByNameFunc: func(ctx context.Context, productName string) (models.Product, error) {
				return models.Product{}, nil
			},
			DeleteProductFromBuyListByIdFunc: func(ctx context.Context, productId string, userid string) error {
				return nil
			},
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
		p:  models.Product{},
		f:  models.FridgeProduct{},
		ur: models.Usertg{},
		ps: models.PeriodStat{},
	}

	// msg := tgbotapi.NewMessage(ChatID, "oaoaao")
	// m, err := bot.BotAPI.Send(msg)
	// if err != nil {
	// 	panic(err)
	// }
	// update, _ := bot.BotAPI.GetUpdates(tgbotapi.UpdateConfig{})

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60

	// updates := bot.BotAPI.GetUpdatesChan(u)
	// var update tgbotapi.Update
	// for update = range updates {
	// 	break
	// }
	// update.CallbackQuery.ID = "404723786"
	// update.CallbackQuery.Data = "deleteProductFromBuyList"

	logger = zap.NewExample()
	defer logger.Sync()
	schb = gocron.NewScheduler(time.Local)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

	msg := tgbotapi.NewMessage(ChatID, "oaoaao")
	if _, err := bot.BotAPI.Send(msg); err != nil {
		t.Error(err)
	}
	bot.HandleCallbacks(context.Background(), &update)
	// bot.BotAPI.MakeRequest(getUpdates)
	// getting too old callback id

	// if _, err := bot.BotAPI.Request(callback); err != nil {
	// 	t.Errorf("callback not handled %s", err)
	// }

	bot.BotAPI.Request(callback)
}
