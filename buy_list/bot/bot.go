package bot

import (
	"buy_list/bot/store"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot() {
	// fmt.Println(os.Getenv("dbuser"))
	// fmt.Println(os.Getenv("tgtoken"))
	// fmt.Println(os.Getenv("dbpassword"))
	token := GetToken()
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
	f := store.FridgeProduct{}

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			//log
			//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.IsCommand() {
				msg.ReplyToMessageID = update.Message.MessageID
				switch update.Message.Command() {
				case "start":
					msg.ReplyMarkup = startKeyboard
					msg.Text = "Привет! Я бот, который может управлять вашими покупками и мониторить срок годности продуктов."
					s.AddUsertg(&store.Usertg{
						Username: update.Message.From.UserName,
					})
				case "cancel":
					GlobalState = StateStart
				default:
					msg.Text = "Неверная команда :("
				}
				SendMessage(bot, &msg)

			} else {
				//reply
				//msg.ReplyToMessageID = update.Message.MessageID

				switch GlobalState {

				//start menu
				case StateStart:
					switch update.Message.Text {
					case startKeyboard.Keyboard[0][0].Text:
						GlobalState = StateAddBuyList
						msg.ReplyMarkup = buylistKeyboard
						SendMessage(bot, &msg)
					case startKeyboard.Keyboard[1][0].Text:
						GlobalState = StateAddFridge
						msg.ReplyMarkup = fridgeKeyboard
						SendMessage(bot, &msg)
					case "open":
						msg.ReplyMarkup = startKeyboard
						SendMessage(bot, &msg)
					case "close":
						msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
						SendMessage(bot, &msg)
					}

				//adding to buy list
				case StateAddBuyList:
					switch update.Message.Text {

					case buylistKeyboard.Keyboard[0][0].Text: //add product
						p.State = 0
						p.UserId = s.GetUseridByUsername(update.Message.From.UserName)
						msg.Text = "Введите название продукта"

					case buylistKeyboard.Keyboard[1][0].Text: //get buy list
						products := s.GetBuyListByUsername(update.Message.From.UserName)
						if len(products) != 0 {
							for i, pr := range products {
								msg.Text = fmt.Sprintf("%d: %s %.2f %s\n", i+1, pr.Name, pr.Weight, pr.BuyDate)
								msg.ReplyMarkup = inlineBuylistKeyboard
								SendMessage(bot, &msg)
							}
							continue
						} else {
							msg.Text = "Список покупок пуст"
							//SendMessage(bot, &msg)
						}

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						GlobalState = StateStart
						msg.ReplyMarkup = startKeyboard

					default:
						switch p.State {

						case 0:
							p.Name = update.Message.Text
							p.ProductId = s.CreateProductByName(p.Name)
							msg.Text = "Введите вес/количество"
							p.State = 1

						case 1:
							if p.Weight, err = strconv.ParseFloat(update.Message.Text, 64); err != nil {
								msg.Text = "Неверный формат"
							} else {
								msg.Text = "Введите время покупки (YYYY-MM-DD HH:MM)"
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
								msg.ReplyMarkup = buylistKeyboard
							}
						}

					}
					SendMessage(bot, &msg)

				case StateAddFridge:
					switch update.Message.Text {

					case fridgeKeyboard.Keyboard[0][0].Text: //add product
						f.State = 0
						f.UserId = s.GetUseridByUsername(update.Message.From.UserName)
						msg.Text = "Введите название продукта"

					case buylistKeyboard.Keyboard[1][0].Text: //get fridge list
						fridgeProducts := s.GetFridgeListByUsername(update.Message.From.UserName)
						if len(fridgeProducts) != 0 {
							for i, pr := range fridgeProducts {
								resText := fmt.Sprintf("%d: %s \n", i+1, pr.Name)
								if f.Opened {
									resText += "Открыт "
								} else {
									resText += "Не вскрыт "
								}
								// NOW add expire and other dates
								expdate, err := time.Parse("2006-01-02", pr.Expire_date)
								if err != nil {
									log.Fatal(err)
								}

								if time.Now().After(expdate) {
									resText += fmt.Sprintf("Просрочен %s ", pr.Expire_date)
								} else {
									resText += fmt.Sprintf("Годен до %s ", pr.Expire_date)
								}

								msg.Text = resText
								msg.ReplyMarkup = inlineBuylistKeyboard
								if _, err := bot.Send(msg); err != nil {
									log.Fatal(err)
								}
							}
						} else {
							msg.Text = "Холодильник пуст"
							SendMessage(bot, &msg)
						}
						continue

					case buylistKeyboard.Keyboard[1][1].Text: //cancel
						GlobalState = StateStart
						msg.ReplyMarkup = startKeyboard

					default:
						switch f.State {
						case 0:
							f.Name = update.Message.Text
							f.ProductId = s.CreateProductByName(f.Name)
							msg.Text = "Укажите срок годности (YYYY-MM-DD)"
							f.State = 1

						case 1:
							ts := update.Message.Text
							if _, err := time.Parse("2006-01-02", ts); err != nil {
								msg.Text = "Неверный формат"
							} else {
								f.Expire_date = ts
								s.AddProductToFridge(&f)
								msg.Text = "Товар добавлен в холодильник"
								msg.ReplyMarkup = fridgeKeyboard
							}

						case 2: //for adding from buy list
							ts := update.Message.Text
							if _, err := time.Parse("2006-01-02", ts); err != nil {
								msg.Text = "Неверный формат"
							} else {
								f.Expire_date = ts
								s.AddProductToFridge(&f)
								pid := s.GetProductIdByName(f.Name)
								s.DeleteProductFromBuyListById(pid)
								msg.Text = "Товар добавлен в холодильник"
								msg.ReplyMarkup = buylistKeyboard
								GlobalState = StateAddBuyList
							}

						}

					}
					SendMessage(bot, &msg)

				}

			}

		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			switch update.CallbackQuery.Data {

			case "deleteProductFromBuyList":
				pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
				pid := s.GetProductIdByName(pname)
				s.DeleteProductFromBuyListById(pid)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					fmt.Sprintf("Продукт '%s' удален из списка покупок", pname))
				SendMessage(bot, &msg)

			case "addToFridgeFromBuyList":
				pname := strings.Fields(update.CallbackQuery.Message.Text)[1]
				pid := s.GetProductIdByName(pname)
				userid := s.GetUseridByUsername(update.CallbackQuery.From.UserName)
				f.Name = pname
				f.UserId = userid
				f.ProductId = pid
				f.Opened = false
				GlobalState = StateAddFridge
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
					"Укажите срок годности (YYYY-MM-DD)")
				SendMessage(bot, &msg)
				f.State = 2
			}

			if _, err := bot.Request(callback); err != nil { //
				log.Fatal(err)
			}
		}
	}
}

func SendMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func GetConn() *store.Store {
	byteToken, err := ioutil.ReadFile("dbdata.json")
	if err != nil {
		log.Fatal("Error during read file: ", err)
	}
	var data store.DBdata
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
	var token store.ApiToken
	err = json.Unmarshal(byteToken, &token)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	return token.Tkn
}
