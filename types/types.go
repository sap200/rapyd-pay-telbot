package types

import (
	"encoding/json"
	"log"
	"sync"
)

var GoRoutingMap = map[string]bool{}

var Wg sync.WaitGroup

type Details struct {
	Name                string         `json:"Name"`
	ChatID              string         `json:"ChatID"`
	Basket              map[string]int `json:"Basket"`
	Amount              float64        `json:"Amount"`
	PayID               string         `json:"PayID"`
	PaymentSubType      string         `json:"PaymentSubType"`
	RedirectURL         string         `json:"RedirectURL"`
	PaymentStatus       string         `json:"PaymentStatus"`
	CardName            string         `json:"CardName"`
	CardNumber          string         `json:"CardNumber"`
	CardExpirationMonth string         `json:"CardExpirationMonth"`
	CardExpirationYear  string         `json:"CardExpirationYear"`
	CardCVV             string         `json:"CardCVV"`
	VPA                 string         `json:"VPA"`
	ShippingAddress     Address        `json:"ShippingAddress"`
}

type Address struct {
	House       string `json:"House"`
	Street      string `json:"Street"`
	State       string `json:"State"`
	City        string `json:"City"`
	Country     string `json:"Country"`
	Postcode    string `json:"Postcode"`
	PhoneNumber string `json:"PhoneNumber"`
}

func (d Details) Marshal() []byte {
	r, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}

	return r
}

func (d *Details) Unmarshal(b []byte) {
	err := json.Unmarshal(b, &d)
	if err != nil {
		log.Println(err)
	}
}

func (d *Details) UpdateAmount() {
	d.Amount = 0
	for k, v := range d.Basket {
		d.Amount += float64(Menu[k] * v)
	}
}
