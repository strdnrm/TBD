package bot

import (
	"buy_list/bot/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {
	token := GetToken()
	fmt.Println(token)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	s := GetConn()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = numericKeyboard
				msg.Text = "Привет! Я бот, который может управлять вашими покупками и мониторить срок годности продуктов."
				s.AddUsertg(store.Usertg{
					Username: update.Message.From.UserName,
				})
			case "test":
				msg.Text = "Test command worked"
			default:
				msg.Text = "Неверная команда :("
			}
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Text {
			case numericKeyboard.Keyboard[0][0].Text:
				msg.Text = "Введите название продукта"
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			case "open":
				msg.ReplyMarkup = numericKeyboard
			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}

			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	}
}

type ApiToken struct {
	Tkn string `json:"token"`
}

type DBdata struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func GetConn() *store.Store {
	byteToken, err := ioutil.ReadFile("dbdata.json")
	if err != nil {
		log.Fatal("Error during read file: ", err)
	}
	var data DBdata
	err = json.Unmarshal(byteToken, &data)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	s := fmt.Sprintf("postgresql://%s:%s@localhost:5433/", data.User, data.Password)
	return store.NewStore(s)
}

func GetToken() string {
	byteToken, err := ioutil.ReadFile("token.json")
	if err != nil {
		log.Fatal("Error during read file: ", err)
	}
	var token ApiToken
	err = json.Unmarshal(byteToken, &token)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return token.Tkn
}
