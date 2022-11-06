package bot

import (
	"buy_list/bot/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

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

	GlobalState := StateStart

	p := store.Product{}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			switch update.Message.Command() {
			case "start":
				msg.ReplyMarkup = numericKeyboard
				msg.Text = "Привет! Я бот, который может управлять вашими покупками и мониторить срок годности продуктов."
				s.AddUsertg(&store.Usertg{
					Username: update.Message.From.UserName,
				})
			case "test":
				msg.Text = "Test command worked"
			case "canse":
				GlobalState = StateStart
			default:
				msg.Text = "Неверная команда :("
			}
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			switch GlobalState {

			//start menu
			case StateStart:
				switch update.Message.Text {
				case numericKeyboard.Keyboard[0][0].Text:
					GlobalState = StateAddBuyList
					p.State = 0
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					p.UserId = s.GetUserid(update.Message.From.UserName)
					// p.BuyListId, p.UserId = s.CreateNewBuyList(update.Message.From.UserName)
					msg.Text = "Введите название продукта"
				case "open":
					msg.ReplyMarkup = numericKeyboard
				case "close":
					msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				}

			//adding to buy list
			case StateAddBuyList:
				switch p.State {
				case 0:
					p.Name = update.Message.Text
					p.ProductId = s.GetProductId(p.Name)
					msg.Text = "Введите вес/количество"
					p.State = 1
				case 1:
					if p.Weight, err = strconv.ParseFloat(update.Message.Text, 64); err != nil {
						msg.Text = "Неверный формат"
					} else {
						msg.Text = "Введите когда время покупки (YYYY-MM-DD HH:MM)"
						p.State = 2
					}
				case 2:
					ts := update.Message.Text + ":00"
					if _, err := time.Parse("2006-01-02 15:04:05.999", ts); err != nil {
						msg.Text = "Неверный формат"
					} else {
						p.BuyDate = ts
						s.AddProductToBuyList(&p)
						msg.Text = "Товар добавлен в список покупок"
						GlobalState = StateStart
						msg.ReplyMarkup = numericKeyboard
					}

				}

			case StateAddFridge:
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
