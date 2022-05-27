package rapyd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/sap200/rapyd-pay-telbot/types"
	"github.com/sap200/rapyd-pay-telbot/utils"
)

// PaymentTypes represents the supported payment with categories
var PaymentTypes = map[string][]string{
	"Credit Card": {"VISA Credit Card", "Amex Credit Card"},
	"Debit Card":  {"Visa Debit Card", "Rupay Debit Card", "Maestro Debit Card"},
	"EWallets":    {"Paytm eWallet", "PhonePe eWallet"},
	"UPI":         {"UPI", "Google Pay via UPI"},
	"Bank":        {"ICICI Bank", "HDFC Bank"},
}

// NameToType matches the payment category to its api request name
var NameToType = map[string]string{
	"VISA Credit Card":   "in_visa_credit_card",
	"Amex Credit Card":   "in_amex_credit_card",
	"Visa Debit Card":    "in_debit_visa_card",
	"Rupay Debit Card":   "in_rupay_debit_card",
	"Maestro Debit Card": "in_maestro_debit_card",
	"Paytm eWallet":      "in_paytm_ewallet",
	"PhonePe eWallet":    "in_phonepe_ewallet",
	"Google Pay via UPI": "in_googlepay_upi_bank",
	"UPI":                "in_upi_bank",
	"ICICI Bank":         "in_icici_bank",
	"HDFC Bank":          "in_hdfc_bank",
}

// GetPrimaryPaymentMethods returns a string of array
// that contains all 5 primary key values in PaymentTypes
func GetPrimaryPaymentMethods() []string {
	res := []string{}
	for k := range PaymentTypes {
		res = append(res, k)
	}

	return res
}

// PerformBasicCheck performs basic check
// for the card fields
func PerformBasicCheck(bot *tgbotapi.BotAPI, chatId int64, arr []string) bool {
	msg := tgbotapi.NewMessage(chatId, "")

	cardNoRe := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`)
	if !cardNoRe.MatchString(arr[0]) {
		msg.Text = "❌  Card Number invalid"
		sendMessage(bot, msg)
		return false
	}

	expMonth, err := strconv.Atoi(arr[2])
	if err != nil {
		msg.Text = "❌  expiration month error"
		sendMessage(bot, msg)
		return false
	}

	if expMonth <= 0 && expMonth >= 12 {
		msg.Text = "❌  expiration month error"
		sendMessage(bot, msg)
		return false
	}

	expYear, err := strconv.Atoi(arr[3])
	if err != nil {
		return false
	}

	year, _, _ := time.Now().Date()
	if expYear < year {
		msg.Text = "❌  expiration Year error"
		sendMessage(bot, msg)
		return false
	}

	cvv, err := strconv.Atoi(arr[4])
	if err != nil {
		msg.Text = "❌  CVV error"
		sendMessage(bot, msg)
		return false
	}

	if cvv <= 99 && cvv >= 1000 {
		msg.Text = "❌   CVV error"
		sendMessage(bot, msg)
		return false
	}

	return true
}

// It uses regex to check VPA format
func PerformBasicCheckUPI(x string) bool {
	upiRE := regexp.MustCompile("^(.+)@(.+)$")
	return upiRE.MatchString(x)
}

// ProcessPayment processes the payment request
// creates a payment request and sends back required field
func ProcessPayment(bot *tgbotapi.BotAPI, chatId int64, dbName string) {
	// access database and get payment sub type
	byteKey := []byte(fmt.Sprintf("%d", chatId))
	d, err := utils.GetDetails(dbName, byteKey)
	if err != nil {
		msg := tgbotapi.NewMessage(chatId, "Details not found")
		sendMessage(bot, msg)
		return
	}
	paymentSubType := d.PaymentSubType

	// if basket is nil

	if d.Basket == nil {
		msg := tgbotapi.NewMessage(chatId, "❌  No items in your Basket. Add items to checkout")
		sendMessage(bot, msg)
		return
	}

	count := 0
	for range d.Basket {
		count++
	}
	if count == 0 {
		msg := tgbotapi.NewMessage(chatId, "❌  No items in your Basket. Add items to checkout")
		sendMessage(bot, msg)
		return
	}

	if strings.Contains(paymentSubType, types.CARD) {
		// if card details are missing then
		if d.CardName == "" || d.CardNumber == "" || d.CardExpirationMonth == "" || d.CardExpirationYear == "" || d.CardCVV == "" {
			msg := tgbotapi.NewMessage(chatId, "")
			msg.Text = "Use /card command and furnish <b><i>Card Number</i></b>, <b><i>CardHolder Name</i></b>, <b><i>Expiration Month</i></b>, <b><i>Expiration Year</i></b>, <b><i>CVV</i></b> in same order and format"
			msg.ParseMode = "html"
			sendMessage(bot, msg)
		} else {
			// proceed to payment
			SendConfirmButton(bot, chatId, d.Amount, dbName)
			// send invoice
		}

	} else if strings.Contains(paymentSubType, types.EWALLET) || strings.Contains(paymentSubType, types.BANK) {

		SendConfirmButton(bot, chatId, d.Amount, dbName)

	} else if strings.Contains(paymentSubType, types.UPI) {
		if d.VPA == "" {
			msg := tgbotapi.NewMessage(chatId, "")
			msg.Text = "Use /vpa <b><i>upiID</i></b> to give upi address you want to pay from"
			msg.ParseMode = "html"
			sendMessage(bot, msg)
		} else {
			SendConfirmButton(bot, chatId, d.Amount, dbName)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, "❌  Payment Type is not defined. Use checkout to proceed.")
		sendMessage(bot, msg)
	}
}

// sendMessage sends a message to telegram bot using telegram bot api
func sendMessage(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

// CreatePaymentRequest creates a payment request
// based on different kind of payment methods.
func CreatePaymentRequest(details types.Details) (string, float64, string, error) {

	if strings.Contains(details.PaymentSubType, types.CARD) {
		body := struct {
			Amount             string                 `json:"amount"`
			Currency           string                 `json:"currency"`
			PaymentMethod      map[string]interface{} `json:"payment_method"`
			CompletePaymentURL string                 `json:"complete_payment_url"`
			ErrorPaymentURL    string                 `json:"error_payment_url"`
		}{
			fmt.Sprintf("%.2f", details.Amount),
			"inr",
			map[string]interface{}{
				"type": NameToType[details.PaymentSubType],
				"fields": map[string]string{
					"number":           details.CardNumber,
					"expiration_month": details.CardExpirationMonth,
					"expiration_year":  details.CardExpirationYear,
					"cvv":              details.CardCVV,
					"name":             details.CardName,
				},
			},

			"https://telegram.me/RapydPBot",
			"https://telegram.me/RapydPBot",
		}
		data, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPost, types.Base_url+"/v1/payments", bytes.NewBuffer(data))
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		c := http.DefaultClient
		signer := NewRapydSigner([]byte(types.Access_key), []byte(types.Secret_key))
		signer.SignRequest(req, data)

		res, err := c.Do(req)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		d, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(d))
		newData := map[string]interface{}{}
		err = json.Unmarshal([]byte(d), &newData)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		payId := newData["data"].(map[string]interface{})["id"]
		amt := newData["data"].(map[string]interface{})["original_amount"]
		redirectURL := newData["data"].(map[string]interface{})["redirect_url"]

		return payId.(string), amt.(float64), redirectURL.(string), nil
	} else if strings.Contains(details.PaymentSubType, types.EWALLET) || strings.Contains(details.PaymentSubType, types.BANK) {
		body := struct {
			Amount             string                 `json:"amount"`
			Currency           string                 `json:"currency"`
			PaymentMethod      map[string]interface{} `json:"payment_method"`
			CompletePaymentURL string                 `json:"complete_payment_url"`
			ErrorPaymentURL    string                 `json:"error_payment_url"`
		}{
			fmt.Sprintf("%.2f", details.Amount),
			"INR",
			map[string]interface{}{
				"type": NameToType[details.PaymentSubType],
			},
			"https://telegram.me/RapydPBot",
			"https://telegram.me/RapydPBot",
		}

		data, _ := json.Marshal(body)
		fmt.Println(string(data))
		req, err := http.NewRequest(http.MethodPost, types.Base_url+"/v1/payments", bytes.NewBuffer(data))
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		c := http.DefaultClient
		signer := NewRapydSigner([]byte(types.Access_key), []byte(types.Secret_key))
		signer.SignRequest(req, data)

		res, err := c.Do(req)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		d, _ := ioutil.ReadAll(res.Body)
		newData := map[string]interface{}{}
		err = json.Unmarshal([]byte(d), &newData)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		fmt.Println(newData)

		payId := newData["data"].(map[string]interface{})["id"]
		amt := newData["data"].(map[string]interface{})["original_amount"]
		redirectURL := newData["data"].(map[string]interface{})["redirect_url"]

		return payId.(string), amt.(float64), redirectURL.(string), nil
	} else if strings.Contains(details.PaymentSubType, types.UPI) {
		body := struct {
			Amount             string                 `json:"amount"`
			Currency           string                 `json:"currency"`
			PaymentMethod      map[string]interface{} `json:"payment_method"`
			CompletePaymentURL string                 `json:"complete_payment_url"`
			ErrorPaymentURL    string                 `json:"error_payment_url"`
		}{
			fmt.Sprintf("%.2f", details.Amount),
			"inr",
			map[string]interface{}{
				"type": NameToType[details.PaymentSubType],
				"fields": map[string]string{
					"vpa": details.VPA,
				},
			},

			"https://telegram.me/RapydPBot",
			"https://telegram.me/RapydPBot",
		}

		data, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPost, types.Base_url+"/v1/payments", bytes.NewBuffer(data))
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		c := http.DefaultClient
		signer := NewRapydSigner([]byte(types.Access_key), []byte(types.Secret_key))
		signer.SignRequest(req, data)

		res, err := c.Do(req)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		d, _ := ioutil.ReadAll(res.Body)
		newData := map[string]interface{}{}
		err = json.Unmarshal([]byte(d), &newData)
		if err != nil {
			log.Println(err)
			return "", 0.0, "", err
		}

		payId := newData["data"].(map[string]interface{})["id"]
		amt := newData["data"].(map[string]interface{})["original_amount"]
		redirectURL := newData["data"].(map[string]interface{})["redirect_url"]

		fmt.Println(string(d))

		return payId.(string), amt.(float64), redirectURL.(string), nil

	} else {
		return "nil", 0.00, "", errors.New("Not a payment type")
	}
}

// Sends a confirmation inline keyboard to telegram
func SendConfirmButton(bot *tgbotapi.BotAPI, chatId int64, originalAmount float64, dbName string) {
	numericKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳   Confirm Payment of ₹ "+fmt.Sprintf("%.3f", originalAmount), types.PAYMENT_CONFIRM+fmt.Sprintf("%d", chatId)),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳  ❌  Reject Payment of ₹ "+fmt.Sprintf("%.3f", originalAmount), types.PAYMENT_REJECT+fmt.Sprintf("%d", chatId)),
		),
	)

	msg := tgbotapi.NewMessage(chatId, "")
	byteKey := []byte(fmt.Sprintf("%d", chatId))
	msg.Text = utils.GetInvoice(dbName, byteKey)
	msg.ReplyMarkup = numericKeyboard
	msg.ParseMode = "html"

	sendMessage(bot, msg)
}

// GetPaymentStatus returns payment status by taking in the payment id as input
func GetPaymentStatus(payID string) string {
	signer := NewRapydSigner([]byte(types.Access_key), []byte(types.Secret_key))
	request, _ := http.NewRequest("GET", types.Base_url+"/v1/payments/"+payID, nil)
	err := signer.SignRequest(request, nil)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	c := http.DefaultClient
	resp, err := c.Do(request)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	x, _ := ioutil.ReadAll(resp.Body)
	result := map[string]interface{}{}
	err = json.Unmarshal(x, &result)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	status := result["data"].(map[string]interface{})["status"]

	return status.(string)

}

// CompletePayment is a dummy function that simulates completion of payment
// for sandbox
func CompletePayment(payID string) string {
	body := map[string]string{
		"token":  payID,
		"param1": "rapyd",
		"param2": "success",
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, types.Base_url+"/v1/payments/completePayment", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return types.PAYMENT_ACTIVE
	}

	c := http.DefaultClient
	signer := NewRapydSigner([]byte(types.Access_key), []byte(types.Secret_key))
	signer.SignRequest(req, data)

	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return types.PAYMENT_ACTIVE
	}

	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
	newData := map[string]interface{}{}
	err = json.Unmarshal([]byte(d), &newData)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	status := newData["data"].(map[string]interface{})["status"]
	return status.(string)
}
