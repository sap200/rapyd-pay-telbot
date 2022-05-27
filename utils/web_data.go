package utils

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sap200/rapyd-pay-telbot/types"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/api/option"
)

func GetAndDelete(dbName string, bot *tgbotapi.BotAPI, chatId string) {

	defer types.Wg.Done()

	ctx := context.Background()
	opt := option.WithCredentialsFile("./types/cred.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	for {
		d, err := client.Doc("webDetail/" + chatId).Get(ctx)
		fmt.Println(d.Data())
		if err == nil {
			_, err = client.Doc("webDetail/" + chatId).Delete(ctx)
			if err != nil {
				fmt.Println(err)
			}

			newData := d.Data()

			fmt.Println(newData)

			err = updateLevelDB(dbName, newData, chatId)
			if err != nil {
				cid, _ := strconv.ParseInt(chatId, 10, 64)
				msg := tgbotapi.NewMessage(cid, "")
				msg.ParseMode = "html"
				msg.Text = "❌   <i>Failed to update details. Please try again.</i>"
				sendMessage(bot, msg)

				// delete go routine
				delete(types.GoRoutingMap, chatId)

				break
			} else {
				// updated details successfully
				// failed to update details
				cid, _ := strconv.ParseInt(chatId, 10, 64)
				msg := tgbotapi.NewMessage(cid, "")
				msg.ParseMode = "html"
				msg.Text = "✅   <i>Details Updated Successfully</i>"
				sendMessage(bot, msg)

				delete(types.GoRoutingMap, chatId)
				break
			}
		}

		time.Sleep(time.Second * 2)

	}

}

func updateLevelDB(dbName string, nd map[string]interface{}, chatId string) error {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get([]byte(chatId), nil)
	if err != nil {
		// data doesn't exists create new detail and put it
		details := types.Details{}
		details.ChatID = chatId
		details.Name = nd["name"].(string)
		details.CardNumber = nd["cardNumber"].(string)
		details.CardName = nd["cardName"].(string)
		details.CardExpirationMonth = nd["expiryMonth"].(string)
		details.CardExpirationYear = nd["expiryYear"].(string)
		details.CardCVV = nd["cvv"].(string)
		details.VPA = nd["vpa"].(string)
		// shipping address
		details.ShippingAddress.PhoneNumber = nd["phone"].(string)
		details.ShippingAddress.House = nd["house"].(string)
		details.ShippingAddress.Street = nd["street"].(string)
		details.ShippingAddress.City = nd["city"].(string)
		details.ShippingAddress.State = nd["state"].(string)
		details.ShippingAddress.Country = nd["country"].(string)
		details.ShippingAddress.Postcode = nd["postcode"].(string)

		// marshall
		val := details.Marshal()
		return check(db.Put([]byte(chatId), []byte(val), nil))

	} else {
		// update existing detail
		var details types.Details
		details.Unmarshal(data)

		// update new details
		details.ChatID = chatId
		details.Name = nd["name"].(string)
		details.CardNumber = nd["cardNumber"].(string)
		details.CardName = nd["cardName"].(string)
		details.CardExpirationMonth = nd["expiryMonth"].(string)
		details.CardExpirationYear = nd["expiryYear"].(string)
		details.CardCVV = nd["cvv"].(string)
		details.VPA = nd["vpa"].(string)
		// shipping address
		details.ShippingAddress.PhoneNumber = nd["phone"].(string)
		details.ShippingAddress.House = nd["house"].(string)
		details.ShippingAddress.Street = nd["street"].(string)
		details.ShippingAddress.City = nd["city"].(string)
		details.ShippingAddress.State = nd["state"].(string)
		details.ShippingAddress.Country = nd["country"].(string)
		details.ShippingAddress.Postcode = nd["postcode"].(string)

		// marshall
		val := details.Marshal()
		return check(db.Put([]byte(chatId), []byte(val), nil))
	}

}

func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
