package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

var schb *gocron.Scheduler //for buy_list alerts
var schf *gocron.Scheduler // for fridge alerts

func InitScheduler(bot *Bot) {
	schb = gocron.NewScheduler(time.Local)
	//for restarting bot
	jbStart, err := schb.Every(1).Seconds().Do(CreateBuyAlerts, bot)
	if err != nil {
		panic(err)
	}
	jbStart.LimitRunsTo(1)

	schb.Every(1).Day().At("00:00").Do(CreateBuyAlerts, bot)
	schb.StartAsync()

	schf = gocron.NewScheduler(time.Local)
	//for restarting bot
	jfStart, err := schf.Every(1).Seconds().Do(CreateExpireAlerts, bot)
	if err != nil {
		panic(err)
	}
	jfStart.LimitRunsTo(1)

	schf.Every(1).Day().At("08:00;18:00").Do(CreateExpireAlerts, bot)
	schf.StartAsync()
}

// for adding new products in buy list during current day
func UpdateBuyListSchedule(bot *Bot) {
	schb.Clear()
	schb.Every(1).Day().At("00:00").Do(CreateBuyAlerts, bot.s, bot)
	CreateBuyAlerts(bot)
	schb.StartAsync()
}

func CreateBuyAlerts(bot *Bot) {
	products, err := bot.s.GetTodayBuyList(context.Background())
	if err != nil {
		logger.Error("Get todays buy list errot", zap.Error(err))
	}
	for _, pr := range products {

		tm, err := time.Parse(time.RFC3339, pr.BuyDate)
		if err != nil {
			panic(err)
		}

		chatid, err := bot.s.GetChatIdByUserId(context.Background(), pr.UserId)
		if err != nil {
			logger.Error("Get chat id error", zap.Error(err))
		}

		text := fmt.Sprintf("Время покупки %s", pr.Name)
		schb.SingletonMode()
		job, err := schb.Every(1).Day().At(tm).Do(SendAlert, bot, chatid, text)
		if err != nil {
			panic(err)
		}
		logger.Info("planned new alert",
			zap.String("Time to run", job.NextRun().String()),
		)
		job.LimitRunsTo(1)
	}
	logger.Info("alerts created",
		zap.Int("jobs count", len(schb.Jobs())),
	)
}

func CreateExpireAlerts(bot *Bot) {
	products, err := bot.s.GetSoonExpireList(context.Background())
	if err != nil {
		logger.Error("Get todays expire list error", zap.Error(err))
	}
	for _, pr := range products {
		tm, err := time.Parse(time.RFC3339, pr.Expire_date)
		if err != nil {
			panic(err)
		}

		chatid, err := bot.s.GetChatIdByUserId(context.Background(), pr.UserId)
		if err != nil {
			logger.Error("Get chat id error", zap.Error(err))
		}

		text := fmt.Sprintf("Скоро истекает срок годности %s %s", pr.Name,
			tm.Format("2006-01-02"))

		SendAlert(bot, chatid, text)
	}
}

func UpdateExpireSchedule(bot *Bot) {
	schf.Clear()
	schf.Every(1).Day().At("08:00;18:00").Do(CreateExpireAlerts, bot)
	schf.StartAsync()
}

func SendAlert(bot *Bot, chat_id int64, text string) {
	msg := tgbotapi.NewMessage(chat_id, text)
	if _, err := bot.BotAPI.Send(msg); err != nil {
		panic(err)
	}
}
