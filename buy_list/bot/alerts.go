package bot

import (
	"buy_list/bot/store"
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var sch *gocron.Scheduler

func InitScheduler(s *store.Store, bot *tgbotapi.BotAPI) {
	// location, err := time.LoadLocation("Europe/Moscow")
	// if err != nil {
	// 	panic(err)
	// }
	sch = gocron.NewScheduler(time.Local)

	//for restarting bot
	fStart, err := sch.Every(1).Seconds().Do(CreateBuyAlerts, s, bot)
	if err != nil {
		panic(err)
	}
	fStart.LimitRunsTo(1)

	sch.Every(1).Day().At("00:00").Do(CreateBuyAlerts, s, bot)
	sch.StartAsync()
}

// for adding new products in buy list during current day
func UpdateBuyListSchedule(s *store.Store, bot *tgbotapi.BotAPI) {
	sch.Clear()
	sch.Every(1).Day().At("00:00").Do(CreateBuyAlerts, s, bot)
	CreateBuyAlerts(s, bot)
	sch.StartAsync()
}

func CreateBuyAlerts(s *store.Store, bot *tgbotapi.BotAPI) {
	products, err := s.GetTodayBuyList(context.Background())
	fmt.Println(len(products))
	if err != nil {
		logger.Error("Get todays buy list errot", zap.Error(err))
	}
	for _, pr := range products {

		tm, err := time.Parse(time.RFC3339, pr.BuyDate)
		if err != nil {
			panic(err)
		}

		chatid, err := s.GetChatIdByUserId(context.Background(), pr.UserId)
		if err != nil {
			logger.Error("Get chat id error", zap.Error(err))
		}

		text := fmt.Sprintf("Время покупки %s", pr.Name)
		sch.SingletonMode()
		job, err := sch.Every(1).Day().At(tm).Do(SendAlert, pr, bot, chatid, text)
		if err != nil {
			panic(err)
		}
		logger.Info("planned new alert",
			zap.String("Time to run", job.NextRun().String()),
		)
		job.LimitRunsTo(1)
	}
	logger.Info("alerts created",
		zap.Int("jobs count", len(sch.Jobs())),
	)
}

func SendAlert(p store.Product, bot *tgbotapi.BotAPI, chat_id int64, text string) {
	msg := tgbotapi.NewMessage(chat_id, text)
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
}
