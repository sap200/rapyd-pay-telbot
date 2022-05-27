package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sap200/rapyd-pay-telbot/rapyd"
	"github.com/sap200/rapyd-pay-telbot/types"
	"github.com/sap200/rapyd-pay-telbot/utils"
)

//var PaymentDetails = map[int64]map[string]string{}

func main() {
	bot, err := tgbotapi.NewBotAPI(types.Telegram_Bot_Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		// Extract the command from the Message.
		if update.Message != nil {
			switch update.Message.Command() {
			case "help":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = types.GetAvailableCommands()
				msg.ParseMode = "html"
				sendMessage(bot, msg)

			case "initstore":
				// fill the map
				chatId := fmt.Sprintf("%d", update.Message.Chat.ID)
				_, ok := types.GoRoutingMap[chatId]
				if !ok {
					types.GoRoutingMap[chatId] = true
					// launch go routine
					types.Wg.Add(1)
					go utils.GetAndDelete(types.DB_Details_Path, bot, chatId)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "✅   <b>Initialized storing...</b>"
					msg.ParseMode = "html"
					sendMessage(bot, msg)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗   <b>Your updations Already initialized...</b>"
					msg.ParseMode = "html"
					sendMessage(bot, msg)
				}

			case "save":
				// fill the map
				chatId := fmt.Sprintf("%d", update.Message.Chat.ID)
				_, ok := types.GoRoutingMap[chatId]
				if !ok {
					types.GoRoutingMap[chatId] = true
					// launch go routine
					types.Wg.Add(1)
					go utils.GetAndDelete(types.DB_Details_Path, bot, chatId)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "✅   <b>Initialized storing...</b>"
					msg.ParseMode = "html"
					sendMessage(bot, msg)
				} else {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗   <b>Your updations Already initialized...</b>"
					msg.ParseMode = "html"
					sendMessage(bot, msg)
				}

			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = types.GetAvailableCommands()
				msg.ParseMode = "html"
				sendMessage(bot, msg)

			case "menu":
				// concrete implementation of menu
				// sends menu to the user as a form of sequence of message
				chatID := update.Message.Chat.ID
				byteKey := []byte(fmt.Sprintf("%d", chatID))
				utils.SetChatID(types.DB_Details_Path, byteKey)
				types.SendMenu(bot, chatID)

			case "summary":
				// send an invoice without checkout option that is supposed to show the summary of basket
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				key := []byte(fmt.Sprintf("%d", update.Message.Chat.ID))
				msg.Text = utils.GetSummary(types.DB_Details_Path, key)
				msg.ParseMode = "html"
				sendMessage(bot, msg)

			case "sethouse":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/sethouse") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/sethouse")+1:], " ")
				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.HOUSE_FIELD, txt)

			case "setstreet":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setstreet") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setstreet")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.STREET_FIELD, txt)

			case "setcity":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setcity") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setcity")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.CITY_FIELD, txt)

			case "setstate":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setstate") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setstate")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.STATE_FIELD, txt)

			case "setcountry":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setcountry") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setcountry")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.COUNTRY_FIELD, txt)

			case "setpostcode":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setpostcode") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setpostcode")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.POSTCODE_FIELD, txt)

			case "setname":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setname") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setname")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.NAME_FIELD, txt)

			case "setphone":
				txt := update.Message.Text
				if len(strings.Trim(txt, " ")) == len("/setphone") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				txt = strings.Trim(txt[len("/setphone")+1:], " ")

				setShippingDetails(bot, update.Message.Chat.ID, types.DB_Details_Path, types.PHONENUMBER_FIELD, txt)

			case "setcard":
				chatId := update.Message.Chat.ID
				if len(strings.Trim(update.Message.Text, " ")) == len("/setcard") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				text := update.Message.Text[len("/setcard")+1:]
				detailArr := strings.Split(text, ", ")

				if len(detailArr) != 5 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				// perform basic checks
				checked := rapyd.PerformBasicCheck(bot, chatId, detailArr)

				if checked {

					byteKey := []byte(fmt.Sprintf("%d", chatId))
					err := utils.SetCardDetails(types.DB_Details_Path, byteKey, detailArr[0], detailArr[1], detailArr[2], detailArr[3], detailArr[4])
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, "❌  Card updation error. Please try again using /card command")
						sendMessage(bot, msg)
					} else {
						msg := tgbotapi.NewMessage(chatId, "✅  Card Set Successfully.")
						sendMessage(bot, msg)
					}
				}

			case "checkout":
				poll := tgbotapi.NewPoll(update.Message.Chat.ID, types.PAYMENT_METHODS+types.PAYMENT_DELIMITER+fmt.Sprintf("%d", update.Message.Chat.ID), rapyd.GetPrimaryPaymentMethods()...)
				sendPoll(bot, poll)

			case "card":
				chatId := update.Message.Chat.ID
				if len(strings.Trim(update.Message.Text, " ")) == len("/card") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}
				text := update.Message.Text[6:]
				detailArr := strings.Split(text, ", ")

				if len(detailArr) != 5 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue
				}

				// perform basic checks
				checked := rapyd.PerformBasicCheck(bot, chatId, detailArr)

				if !checked {
					rapyd.ProcessPayment(bot, chatId, types.DB_Details_Path)

				}

				if checked {

					byteKey := []byte(fmt.Sprintf("%d", chatId))
					err := utils.SetCardDetails(types.DB_Details_Path, byteKey, detailArr[0], detailArr[1], detailArr[2], detailArr[3], detailArr[4])
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, "❌  Card updation error. Please try again using /card command")
						sendMessage(bot, msg)
					} else {
						msg := tgbotapi.NewMessage(chatId, "✅  Card Set Successfully.")
						sendMessage(bot, msg)
						rapyd.ProcessPayment(bot, chatId, types.DB_Details_Path)

					}

				}

			case "vpa":
				chatId := update.Message.Chat.ID
				if len(strings.Trim(update.Message.Text, " ")) == len("/vpa") {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue

				}
				vpa := update.Message.Text[5:]
				vpa = strings.Trim(vpa, " ")

				check := rapyd.PerformBasicCheckUPI(vpa)
				if !check {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❌   Invalid VPA ! Try again"
					sendMessage(bot, msg)
				}

				if check {
					// set VPA
					byteKey := []byte(fmt.Sprintf("%d", chatId))
					err := utils.SetVPA(types.DB_Details_Path, byteKey, vpa)
					// process payment
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, "❌  VPA updation error. Use /vpa command")
						sendMessage(bot, msg)
					} else {
						msg := tgbotapi.NewMessage(chatId, "✅  VPA Set Successfully.")
						sendMessage(bot, msg)
						rapyd.ProcessPayment(bot, chatId, types.DB_Details_Path)
					}
				}

			case "setvpa":
				chatId := update.Message.Chat.ID
				if len(strings.Trim(update.Message.Text, " ")) == len("/setvpa") {

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❗  Please provide arguments. Check help command for more."
					sendMessage(bot, msg)
					continue

				}
				vpa := update.Message.Text[len("/setvpa")+1:]
				vpa = strings.Trim(vpa, " ")
				check := rapyd.PerformBasicCheckUPI(vpa)
				if !check {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = "❌   Invalid VPA ! Try again"
					sendMessage(bot, msg)
				}

				if check {
					// set VPA
					byteKey := []byte(fmt.Sprintf("%d", chatId))
					err := utils.SetVPA(types.DB_Details_Path, byteKey, vpa)
					// process payment
					if err != nil {
						msg := tgbotapi.NewMessage(chatId, "❌  VPA updation error. Use /vpa command")
						sendMessage(bot, msg)
					} else {
						msg := tgbotapi.NewMessage(chatId, "✅  VPA Set Successfully.")
						sendMessage(bot, msg)
					}
				}

			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.Text = "I don't know that command"
				sendMessage(bot, msg)
			}
		}

		if update.Poll != nil {
			arr := strings.Split(update.Poll.Question, "-")
			qns := arr[0]
			chatId, _ := strconv.ParseInt(arr[1], 10, 64)
			switch qns {
			case types.PAYMENT_METHODS:
				options := update.Poll.Options
				paymentType := ""
				for _, opt := range options {
					if opt.VoterCount > 0 {
						paymentType = opt.Text
					}
				}

				// get the sub Payment methods and send another poll
				poll := tgbotapi.NewPoll(chatId, types.PAYMENT_CATEGORY+types.PAYMENT_DELIMITER+fmt.Sprintf("%d", chatId), rapyd.PaymentTypes[paymentType]...)
				sendPoll(bot, poll)

			case types.PAYMENT_CATEGORY:
				options := update.Poll.Options
				paymentSubType := ""
				for _, opt := range options {
					if opt.VoterCount > 0 {
						paymentSubType = opt.Text
					}
				}

				//PaymentDetails[chatId]["paymentSubType"] = paymentSubType
				// update the Payment type to the database
				byteKey := []byte(fmt.Sprintf("%d", chatId))
				err := utils.SetPaymentSubTypeAndStatus(types.DB_Details_Path, byteKey, paymentSubType, types.PAYMENT_ACTIVE)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.Text = err.Error()
				} else {
					fmt.Println("----------------Displaying Database here--------------------")
					utils.DisplayDBItem(types.DB_Details_Path, byteKey)
					rapyd.ProcessPayment(bot, chatId, types.DB_Details_Path)
				}
			}
		}

		if update.CallbackQuery != nil {
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			// when a request is received and whatever in data if it contains add-, then we need to
			// add to the users basket.
			data := update.CallbackQuery.Data
			chatID := update.CallbackQuery.Message.Chat.ID
			if strings.Contains(data, types.ADD_TO_BASKET) {
				// its a add command
				// split the data and extract itemName
				itemName := strings.Split(data, "-")[1]
				byteChatID := []byte(fmt.Sprintf("%d", chatID))
				// add to basket
				err := utils.AddToBasket(types.DB_Details_Path, byteChatID, itemName)
				fmt.Println(err)
				utils.DisplayDBItem(types.DB_Details_Path, byteChatID)
				if err == nil {
					msg := tgbotapi.NewMessage(chatID, "")
					msg.Text = "🟢  Added 1 " + itemName + " to basket."
					sendMessage(bot, msg)
				} else {
					msg := tgbotapi.NewMessage(chatID, "")
					msg.Text = "❌  There's some error in adding item to basket"
					sendMessage(bot, msg)
				}

			} else if strings.Contains(data, types.REMOVE_FROM_BASKET) {
				// its a remove command
				// split the data and extract itemName
				itemName := strings.Split(data, "-")[1]
				byteChatID := []byte(fmt.Sprintf("%d", chatID))

				err := utils.RemoveFromBasket(types.DB_Details_Path, byteChatID, itemName)
				msg := tgbotapi.NewMessage(chatID, "")

				if err != nil {
					msg.Text = err.Error()
				} else {
					msg.Text = "🔴  Removed from basket"
				}

				sendMessage(bot, msg)

			} else if strings.Contains(data, types.PAYMENT_CONFIRM) {

				byteChatID := []byte(fmt.Sprintf("%d", chatID))
				// get Details
				details, _ := utils.GetDetails(types.DB_Details_Path, byteChatID)
				paymentSubType := details.PaymentSubType

				if strings.Contains(paymentSubType, types.EWALLET) || strings.Contains(paymentSubType, types.UPI) || strings.Contains(paymentSubType, types.BANK) {

					// create payment request
					payID, amt, redirectURL, err := rapyd.CreatePaymentRequest(details)
					if err != nil {
						msg := tgbotapi.NewMessage(chatID, "❌   "+err.Error())
						sendMessage(bot, msg)
						log.Println(err)
						continue
					}

					if strings.Contains(paymentSubType, types.BANK) {
						redirectURL = fmt.Sprintf("https://dashboard.rapyd.net/complete-bank-payment?token=%s&complete_payment_url=%s&error_payment_url=%s", payID, "https://telegram.me/RapydPBot", "https://telegram.me/RapydPBot")
					}

					// set payid, redirect url
					err = utils.SetPayIDRedirectURLAndStatus(types.DB_Details_Path, byteChatID, payID, redirectURL, "")
					if err != nil {
						msg := tgbotapi.NewMessage(chatID, "❌  Error saving values to database.")
						sendMessage(bot, msg)
						continue
					}
					// for eWallet
					// here send a payment redirect link
					// send edit config message with redirect URL as a button

					//amnt := fmt.Sprintf("%.2f", amt)
					newRedirectURL := strings.Split(redirectURL, "&")[0]
					msgLink := fmt.Sprintf(`https://api.telegram.org/bot%s/editMessageReplyMarkup?chat_id=%d&message_id=%d&reply_markup={"inline_keyboard":[[{"text":"✔ Proceed to Pay ₹ %.2f","web_app":{"url":"%s"}}]]}`, types.Telegram_Bot_Token, chatID, update.CallbackQuery.Message.MessageID, amt, newRedirectURL)
					fmt.Println(msgLink)
					res, err := http.Get(msgLink)
					if err != nil {
						msg := tgbotapi.NewMessage(chatID, "❌  Error Editing Message.")
						sendMessage(bot, msg)
						continue
					}
					// no need of message because its already done
					dat, _ := ioutil.ReadAll(res.Body)
					fmt.Println(string(dat))

					// msg1 := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
					// msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteChatID)
					// msg1.ParseMode = "html"
					// numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					// 	tgbotapi.NewInlineKeyboardRow(
					// 		tgbotapi.NewInlineKeyboardButtonURL("✔ Proceed to Pay ₹ "+amnt, redirectURL),
					// 	),
					// )
					// msg1.ReplyMarkup = &numericKeyboard

					// if _, err := bot.Send(msg1); err != nil {
					// 	panic(err)
					// }

					counter := 0

					types.Wg.Add(1)
					go func() {
						defer types.Wg.Done()
						for {
							status := rapyd.GetPaymentStatus(payID)
							// this goes inside go routine with a for loop
							if status == types.PAYMENT_CLOSED {
								// set the PayID and redirectURL and update the payment status
								_ = utils.SetPayIDRedirectURLAndStatus(types.DB_Details_Path, byteChatID, "", "", status)

								// edit the payment detail
								amt := fmt.Sprintf("%.2f", details.Amount)
								msg1 := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
								msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteChatID)
								msg1.ParseMode = "html"
								numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
									tgbotapi.NewInlineKeyboardRow(
										tgbotapi.NewInlineKeyboardButtonData("✅   Payment of ₹ "+amt+" Successful", types.GIBBERISH+"-"+fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)),
									),
								)
								msg1.ReplyMarkup = &numericKeyboard

								if _, err := bot.Send(msg1); err != nil {
									log.Println(err)
								}

								// send a message that order is on the way....

								// checked out and hence move to other database
								oid, _ := utils.MoveFromDB(types.DB_Details_Path, types.DB_Orders_Path, byteChatID)

								m := tgbotapi.NewMessage(chatID, "")
								m.Text = "<i>Your order <pre>" + oid + "</pre> is not on the way. This is a dummy app and hence do not expect your order and be hungry.</i>  😉"
								m.ParseMode = "html"
								sendMessage(bot, m)
								// -------------debug---------------
								// print orders db
								utils.IterateOverDB(types.DB_Orders_Path)
								break

							} else if status == types.PAYMENT_ERROR {
								if counter == 0 {
									msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌ Payment Failed ! Try again")
									sendMessage(bot, msg)
									counter++
								}
							}

							time.Sleep(time.Second * 2)
						}

					}()

				} else {

					// For Card Payments
					// ----------------------------------------------------------------------------------------------------------------------------------------
					// creae a payment request
					// create a map with details for payment request
					payID, amt, redirectURL, err := rapyd.CreatePaymentRequest(details)
					if err != nil {
						msg := tgbotapi.NewMessage(chatID, "❌   "+err.Error())
						sendMessage(bot, msg)
						log.Println(err)
						continue
					}

					if redirectURL == "" {
						// sleep for 2 secs for the payment to get updated
						time.Sleep(time.Second * 4)
						// try retrieving the payment status using paymentID
						status := rapyd.GetPaymentStatus(payID)

						// if result == CLO then do this
						if status == types.PAYMENT_CLOSED {

							// set the PayID and redirectURL and update the payment status
							_ = utils.SetPayIDRedirectURLAndStatus(types.DB_Details_Path, byteChatID, payID, redirectURL, status)

							// edit the payment detail
							amt := fmt.Sprintf("%.2f", amt)
							msg1 := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
							msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteChatID)
							msg1.ParseMode = "html"
							numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
								tgbotapi.NewInlineKeyboardRow(
									tgbotapi.NewInlineKeyboardButtonData("✅   Payment of ₹ "+amt+" Successful", types.GIBBERISH+"-"+fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)),
								),
							)
							msg1.ReplyMarkup = &numericKeyboard

							if _, err := bot.Send(msg1); err != nil {
								log.Println(err)
							}

							// checked out and hence move to other database
							oid, _ := utils.MoveFromDB(types.DB_Details_Path, types.DB_Orders_Path, byteChatID)
							m := tgbotapi.NewMessage(chatID, "")
							m.Text = "<i>Your order <pre>" + oid + "</pre> is not on the way. This is a dummy app and hence do not expect your order and be hungry.</i>  😉"
							m.ParseMode = "html"
							sendMessage(bot, m)

							// -------------debug---------------
							// print orders db
							utils.IterateOverDB(types.DB_Orders_Path)

						} else {
							// else the text is different and its failure and asks to do another payment.
							msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌   Payment Failed ! Try again")
							if _, err := bot.Send(msg); err != nil {
								log.Println(err)
							}

						}
					} else {
						newRedirectURL := strings.Split(redirectURL, "&")[0]
						msgLink := fmt.Sprintf(`https://api.telegram.org/bot%s/editMessageReplyMarkup?chat_id=%d&message_id=%d&reply_markup={"inline_keyboard":[[{"text":"✔ Press for 3DS verification","web_app":{"url":"%s"}}]]}`, types.Telegram_Bot_Token, chatID, update.CallbackQuery.Message.MessageID, newRedirectURL)
						fmt.Println(msgLink)
						res, err := http.Get(msgLink)
						if err != nil {
							msg := tgbotapi.NewMessage(chatID, "❌  Error Editing Message.")
							sendMessage(bot, msg)
							continue
						}
						// no need of message because its already done
						dat, _ := ioutil.ReadAll(res.Body)
						fmt.Println(string(dat))
						// numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
						// 	tgbotapi.NewInlineKeyboardRow(
						// 		tgbotapi.NewInlineKeyboardButtonURL("✔ Press for 3DS verification", redirectURL),
						// 	),
						// )
						// msg1 := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "", numericKeyboard)
						// msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteChatID)
						// msg1.ParseMode = "html"
						// if _, err := bot.Send(msg1); err != nil {
						// 	panic(err)
						// }

						time.Sleep(time.Second * 5)
						counter := 0
						types.Wg.Add(1)
						go func() {
							defer types.Wg.Done()
							for {
								status := rapyd.GetPaymentStatus(payID)
								//fmt.Println("----------------------Inside go routine---------------------")
								fmt.Println(status)
								if status == types.PAYMENT_CLOSED {
									// set the PayID and redirectURL and update the payment status
									_ = utils.SetPayIDRedirectURLAndStatus(types.DB_Details_Path, byteChatID, payID, redirectURL, status)

									amnt := fmt.Sprintf("%.2f", amt)
									msg1 := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "Payment of ₹ "+amnt+" Successful")
									msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteChatID)
									msg1.ParseMode = "html"
									keyb := tgbotapi.NewInlineKeyboardMarkup(
										tgbotapi.NewInlineKeyboardRow(
											tgbotapi.NewInlineKeyboardButtonData("✅   Payment of ₹ "+amnt+" Successful", types.GIBBERISH+"-"+fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)),
										),
									)
									msg1.ReplyMarkup = &keyb
									if _, err := bot.Send(msg1); err != nil {
										log.Println(err)
									}

									// checked out and hence move to other database
									oid, _ := utils.MoveFromDB(types.DB_Details_Path, types.DB_Orders_Path, byteChatID)
									m := tgbotapi.NewMessage(chatID, "")
									m.Text = "<i>Your order <pre>" + oid + "</pre> is not on the way. This is a dummy app and hence do not expect your order and be hungry.</i>  😉"
									m.ParseMode = "html"
									sendMessage(bot, m)

									// -------------debug---------------
									// print orders db
									utils.IterateOverDB(types.DB_Orders_Path)

									break
								} else if status == types.PAYMENT_ERROR {
									if counter == 0 {
										msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "❌   Payment Failed ! Try again")
										sendMessage(bot, msg)
										counter++
									}
								}

								time.Sleep(time.Second * 2)
							}

						}()
					}
				}

			} else if strings.Contains(data, types.PAYMENT_REJECT) {

				// payment was rejected, hence do nothing but edit the button to payment rejected
				// creating a message and setting gibberish so not to hold as a callback query
				msg1 := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
				byteKey := []byte(fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID))
				msg1.Text = utils.GetInvoice(types.DB_Details_Path, byteKey)
				msg1.ParseMode = "html"
				numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("❌  Payment Cancelled", types.GIBBERISH+"-"+fmt.Sprintf("%d", update.CallbackQuery.Message.Chat.ID)),
					),
				)
				msg1.ReplyMarkup = &numericKeyboard
				if _, err := bot.Send(msg1); err != nil {
					log.Println(err)
				}

			}

		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func sendPoll(bot *tgbotapi.BotAPI, poll tgbotapi.SendPollConfig) {
	if _, err := bot.Send(poll); err != nil {
		log.Println(err)
	}
}

func setShippingDetails(bot *tgbotapi.BotAPI, chatID int64, dbName string, field string, fieldValue string) {
	msg := tgbotapi.NewMessage(chatID, "")
	key := []byte(fmt.Sprintf("%d", chatID))
	err := utils.SetShippingAddress(dbName, key, field, fieldValue)
	if err != nil {
		msg.Text = err.Error()
	} else {
		msg.Text = strings.Title(field) + " successfully set"
	}

	sendMessage(bot, msg)
}
