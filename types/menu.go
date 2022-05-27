package types

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nfnt/resize"
)

var Menu = map[string]int{
	"Boneless Strips Bucket 4pcs":  299,
	"Drumstisck Bucket 4pcs":       424,
	"Grill Drumstick 4pcs":         294,
	"Grill Drumstick 8pcs":         578,
	"Fiery Couples Delight Bucket": 308,
	"Chicken Nuggets 3pcs":         74,
	"Chicken Spring Roll 4pcs":     94,
	"Super Veggie Wrap":            74,
	"Grill Chicken Burger":         209,
	"Super Twin Burger":            279,
	"Classic Chicken Burger":       224,
	"DB Grilled Burger":            264,
	"ME Fried Chicken Burger":      224,
	"Mint Mojito":                  64,
	"Green Apple Mojito":           64,
	"Chocolate Frappe":             79,
	"Oreo Cookie Frappe":           104,
	"French Fries":                 49,
	"Tandoori Chicken Sandwich":    109,
}

var ItemDescription = map[string]string{
	"Boneless Strips Bucket 4pcs":  "Banquet of r4 100% chicken breast Mini Fillets. A fuss-free dipper's delight",
	"Drumstisck Bucket 4pcs":       "A drumstick is the lower part of the leg of a bird such as a chicken which is cooked and eaten. A bucket with 4pcs.",
	"Grill Drumstick 4pcs":         "Grill on direct heat for about 30 minutes, depending on the size of the drumsticks. Flip every 10 minutes or so. Cook until the internal temperature of the chicken reaches at least 165° Fahrenheit. Remove from grill and allow the chicken to rest for at least five minutes. A bucket with 4pcs.",
	"Grill Drumstick 8pcs":         "Grill on direct heat for about 30 minutes, depending on the size of the drumsticks. Flip every 10 minutes or so. Cook until the internal temperature of the chicken reaches at least 165° Fahrenheit. Remove from grill and allow the chicken to rest for at least five minutes. A bucket with 8pcs.",
	"Fiery Couples Delight Bucket": "Normally served as an appetizer, Couple's Delight contains thin slices of roast chicken mostly taken from the leg. A bucket with 4pcs.",
	"Chicken Nuggets 3pcs":         "A chicken nugget is a food product consisting of a small piece of deboned chicken meat that is breaded or battered, then deep-fried or baked. Invented in the 1950s.",
	"Chicken Spring Roll 4pcs":     "A spring roll is a Chinese food consisting of a small roll of thin pastry filled with vegetables and sometimes chicken meat, and then fried.",
	"Super Veggie Wrap":            "A wrap is a food dish made with a soft flatbread rolled around a filling. The usual flatbreads are wheat tortillas, lavash, or pita; the filling may include shredded lettuce, diced tomato or pico de gallo, guacamole, sauteed mushrooms, bacon, grilled onions, cheese, and a sauce, such as ranch or honey mustard.",
	"Grill Chicken Burger":         "A Grill Chicke Burger consists of boneless, skinless chicken breast or thigh grilled, served between slices of bread, on a bun,",
	"Super Twin Burger":            "A mix of chicken Fillet, minced chicken patty deep fried, dressed with tomato slice, cheese and bun.",
	"Classic Chicken Burger":       "A juicy lightly breaded crispy chicken breast with crunchy lettuce, tomato, mayonnaise, and the perfect pickles all on a toasted bun.",
	"DB Grilled Burger":            "2 pieces Boneless marinated chicken fillet grilled, dressed with cheese, lettuce, Tandoori mayo and papa bun",
	"ME Fried Chicken Burger":      "ME Fried chicken Burger is a burger that typically consists of deep fried boneless, skinless chicken breast or thigh served between slices of bread, on a bun.",
	"Mint Mojito":                  "The cocktail often consists of five ingredients: white rum, sugar, lime juice, soda water, and mint. Its combination of sweetness, citrus, and herbaceous mint flavors is intended to complement the rum.",
	"Green Apple Mojito":           "Green Apple Mojito is prepared with white rum, green apple syrup, lemon juice, soda and mint leaves.",
	"Chocolate Frappe":             "chocolate frappe will result in a milkshake: ice cream that's blended with syrup and sometimes milk which you can sip through a straw.",
	"Oreo Cookie Frappe":           "Oreo Frappe is a chocolatey cold coffee packed with the goodness and crunchiness of Oreo cookies. It is a heavenly frappe recipe using three different forms of chocolate – coffee, chocolate syrup, and Oreo cookies.",
	"French Fries":                 "A cut of Julienne potatoes deep fried, salted and served with tomato sauce.",
	"Tandoori Chicken Sandwich":    "Tandoori Chicken sandwich is made with layers of tandoori chicken, juicy tomatoes, hard boiled eggs, crispy lettuce, cheddar cheese and a spicy mayonnaise dressing",
}

// SendMenu sends multiple pictures to the telegram channel
// This pictures contains description of product with its price tag.
func SendMenu(bot *tgbotapi.BotAPI, chatID int64) {
	for item := range Menu {

		keyb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🟢  Add", ADD_TO_BASKET+item),
				tgbotapi.NewInlineKeyboardButtonData("🔴  Remove", REMOVE_FROM_BASKET+item),
			),

			// tgbotapi.NewInlineKeyboardRow(
			// 	tgbotapi.NewInlineKeyboardButtonData("🔴  Remove from basket", REMOVE_FROM_BASKET+item),
			// ),
		)

		m := tgbotapi.NewPhoto(chatID, getPhoto(resolvePath(item)))
		m.ParseMode = "html"
		m.Caption = getMessage(item)
		m.ReplyMarkup = keyb

		if _, err := bot.Send(m); err != nil {
			log.Println(err)
		}
	}
}

func resolvePath(item string) string {
	return fmt.Sprintf("./assets/%s.jpg", item)
}

func getMessage(item string) string {
	str := "<u><b>" + item + "</b></u>\n\n"
	str += "<i>" + ItemDescription[item] + "</i>\n\n"
	str += "<b> Price </b>: ₹ " + strconv.Itoa(Menu[item])

	return str
}

func getPhoto(path string) tgbotapi.FileBytes {
	photoBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	image, _, _ := image.Decode(bytes.NewReader(photoBytes))
	newImage := resize.Resize(300, 250, image, resize.Lanczos3)
	buffer := new(bytes.Buffer)
	_ = jpeg.Encode(buffer, newImage, nil)

	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: buffer.Bytes(),
	}
	photoFileBytes.Name = "Burger"

	return photoFileBytes
}
