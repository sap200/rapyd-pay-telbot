package utils

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/sap200/rapyd-pay-telbot/types"
	"github.com/syndtr/goleveldb/leveldb"
)

func GetInvoice(dbName string, key []byte) string {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		return "❌ No Details found about you. Order something and set fields and view the summary again."
	} else {
		str := "🧾  <u><b>INVOICE</b></u>\n\n"

		var d types.Details
		d.Unmarshal(data)

		str += "<u><b>Merchant</b></u>: <i>" + "Daily Belly   🍔" + "</i>\n"
		str += "<u><b>Name</b></u> : <i>" + d.Name + "</i>\n\n"

		if d.Basket != nil {
			count := 0
			b := new(bytes.Buffer)
			writer := tabwriter.NewWriter(b, 2, 8, 3, ' ', 0)
			for range d.Basket {
				count++
			}
			if count != 0 {
				fmt.Fprintln(writer, "Items\tQuantity\tPrice\t")
				fmt.Fprintln(writer, "-----\t---------\t-----\t")
			}
			for k, v := range d.Basket {
				fmt.Fprintln(writer, k+"\t"+fmt.Sprintf("%d", v)+"\t"+fmt.Sprintf("%d", types.Menu[k])+"\t")
			}
			writer.Flush()
			str += "<pre>" + b.String() + "</pre>\n"
		}

		str += "<b>Total</b>: " + fmt.Sprintf("%.2f", d.Amount) + "\n\n"
		str += "<b>shipping address</b> - <i>" + d.ShippingAddress.House + ", " + d.ShippingAddress.Street + ", " + d.ShippingAddress.City + ", " + d.ShippingAddress.State + ", " + d.ShippingAddress.Country + " - " + d.ShippingAddress.Postcode + "</i>\n"
		str += "<b>Phone</b> - <i>" + d.ShippingAddress.PhoneNumber + "</i>\n"

		if strings.Contains(d.PaymentSubType, types.CARD) {
			str += "<b>Payment Method</b>- <i>" + d.PaymentSubType + "  " + d.CardNumber[len(d.CardNumber)-4:len(d.CardNumber)] + "</i>\n\n"
		} else if strings.Contains(d.PaymentSubType, types.EWALLET) || strings.Contains(d.PaymentSubType, types.BANK) {
			str += "<b>Payment Method</b>- <i>" + d.PaymentSubType + "</i>\n\n"
		} else {
			str += "<b>Payment Method</b>- <i>" + d.PaymentSubType + "  " + d.VPA + "</i>\n\n"
		}

		str += "<i><b>Thankyou</b></i>"

		return str
	}
}

func GetSummary(dbName string, key []byte) string {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		return "❌  No Details found about you. Order something and set fields and view the summary again."
	} else {
		str := "📜  <u><b>Summary</b></u>\n\n"

		var d types.Details
		d.Unmarshal(data)

		str += "<u><b>Name</b></u> : <i>" + d.Name + "</i>\n\n"

		if d.Basket != nil {
			count := 0
			b := new(bytes.Buffer)
			writer := tabwriter.NewWriter(b, 2, 8, 3, ' ', 0)
			for range d.Basket {
				count++
			}
			if count != 0 {
				fmt.Fprintln(writer, "Items\tQuantity\tPrice\t")
				fmt.Fprintln(writer, "-----\t--------\t-----\t")
			}
			for k, v := range d.Basket {
				fmt.Fprintln(writer, k+"\t"+fmt.Sprintf("%d", v)+"\t"+fmt.Sprintf("%d", types.Menu[k])+"\t")
			}
			writer.Flush()
			str += "<pre>" + b.String() + "</pre>\n"
		}

		str += "<b>Total</b>: " + fmt.Sprintf("%.2f", d.Amount) + "\n\n"
		str += "<b>shipping address</b> - <i>" + d.ShippingAddress.House + ", " + d.ShippingAddress.Street + ", " + d.ShippingAddress.City + ", " + d.ShippingAddress.State + ", " + d.ShippingAddress.Country + " - " + d.ShippingAddress.Postcode + "</i>\n"
		str += "<b>Phone</b> - <i>" + d.ShippingAddress.PhoneNumber + "</i>\n"

		return str
	}

}
