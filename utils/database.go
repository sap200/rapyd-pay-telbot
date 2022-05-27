package utils

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sap200/rapyd-pay-telbot/types"
	"github.com/syndtr/goleveldb/leveldb"
)

func AddToBasket(dbName string, key []byte, item string) error {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		// key doesnt exists and hence create one
		details := types.Details{}
		if details.ChatID == "" {
			details.ChatID = string(key)
		}
		details.Basket = map[string]int{}
		details.Basket[item] = 1
		details.UpdateAmount()
		val := details.Marshal()

		return check(db.Put(key, val, nil))
	} else {
		// it already exists
		// and hence unmarshal it
		var details types.Details
		details.Unmarshal(data)
		if details.Basket != nil {
			_, ok := details.Basket[item]
			if ok {
				// then key already exists
				details.Basket[item]++
			} else {
				details.Basket[item] = 1
			}

			details.UpdateAmount()
			// marshal details and put back to db
			toPut := details.Marshal()
			return check(db.Put(key, toPut, nil))
		} else {
			details.Basket = map[string]int{}
			if details.ChatID == "" {
				details.ChatID = string(key)
			}
			details.Basket[item] = 1
			details.UpdateAmount()
			val := details.Marshal()
			return check(db.Put(key, val, nil))
		}

	}
}

func RemoveFromBasket(dbName string, key []byte, itemName string) error {
	// open the database
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// check if the key exists
	data, err := db.Get(key, nil)
	if err != nil {
		return errors.New("❗ key not found")
	}

	// if key exists and basket is not nil
	var d types.Details
	d.Unmarshal(data)

	if d.Basket == nil {
		return errors.New("❗  empty basket")
	}

	_, ok := d.Basket[itemName]
	if !ok {
		return errors.New("❗  no such item in basket")
	}

	d.Basket[itemName]--
	if d.Basket[itemName] == 0 {
		delete(d.Basket, itemName)
	}

	d.UpdateAmount()

	val := d.Marshal()
	return check(db.Put(key, val, nil))
}

// calculate basket price and update
func check(err error) error {
	return err
}

func DisplayDBItem(dbName string, key []byte) {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)

	if err != nil {
		fmt.Println("Database Item not found")
	} else {
		fmt.Println(string(data))
	}

}

func SetShippingAddress(dbName string, key []byte, field string, fieldValue string) error {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		// the value doesn't already exists
		details := types.Details{}
		if details.ChatID == "" {
			details.ChatID = string(key)
		}

		switch field {
		case types.HOUSE_FIELD:
			details.ShippingAddress.House = fieldValue
		case types.STREET_FIELD:
			details.ShippingAddress.Street = fieldValue
		case types.CITY_FIELD:
			details.ShippingAddress.City = fieldValue
		case types.COUNTRY_FIELD:
			details.ShippingAddress.Country = fieldValue
		case types.STATE_FIELD:
			details.ShippingAddress.State = fieldValue
		case types.POSTCODE_FIELD:
			details.ShippingAddress.Postcode = fieldValue
		case types.PHONENUMBER_FIELD:
			details.ShippingAddress.PhoneNumber = fieldValue
		case types.NAME_FIELD:
			details.Name = fieldValue
		}

		val := details.Marshal()

		return check(db.Put(key, val, nil))
	} else {
		// unmarshall
		var details types.Details
		details.Unmarshal(data)

		switch field {
		case types.HOUSE_FIELD:
			details.ShippingAddress.House = fieldValue
		case types.STREET_FIELD:
			details.ShippingAddress.Street = fieldValue
		case types.CITY_FIELD:
			details.ShippingAddress.City = fieldValue
		case types.COUNTRY_FIELD:
			details.ShippingAddress.Country = fieldValue
		case types.STATE_FIELD:
			details.ShippingAddress.State = fieldValue
		case types.POSTCODE_FIELD:
			details.ShippingAddress.Postcode = fieldValue
		case types.PHONENUMBER_FIELD:
			details.ShippingAddress.PhoneNumber = fieldValue
		case types.NAME_FIELD:
			details.Name = fieldValue
		}
		val := details.Marshal()

		return check(db.Put(key, val, nil))
	}
}

func SetChatID(dbName string, key []byte) error {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// fetch from database
	data, err := db.Get(key, nil)
	if err != nil {
		// no key exists
		details := types.Details{}
		details.ChatID = string(key)
		val := details.Marshal()
		return check(db.Put(key, val, nil))
	} else {
		var d types.Details
		d.Unmarshal(data)
		if d.ChatID == "" {
			d.ChatID = string(key)
		}

		val := d.Marshal()
		return check(db.Put(key, val, nil))
	}
}

func SetPayIDRedirectURLAndStatus(dbName string, key []byte, payID string, redirectURL string, status string) error {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// fetch data from database
	data, _ := db.Get(key, nil)
	var d types.Details
	d.Unmarshal(data)

	if payID != "" {
		d.PayID = payID
	}
	if redirectURL != "" {
		d.RedirectURL = redirectURL
	}

	if status != "" {
		d.PaymentStatus = status
	}

	// put it into database
	val := d.Marshal()
	return check(db.Put(key, val, nil))

}

func SetPaymentSubTypeAndStatus(dbName string, key []byte, fieldSubType string, fieldStatus string) error {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// fetch data from database
	data, err := db.Get(key, nil)
	if err != nil {
		details := types.Details{}
		details.PaymentSubType = fieldSubType
		details.PaymentStatus = fieldStatus
		val := details.Marshal()
		return check(db.Put(key, val, nil))
	} else {
		var d types.Details
		d.Unmarshal(data)
		d.PaymentSubType = fieldSubType
		d.PaymentStatus = fieldStatus
		val := d.Marshal()
		return check(db.Put(key, val, nil))
	}

}

func SetVPA(dbName string, key []byte, vpa string) error {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// fetch data from database
	data, err := db.Get(key, nil)
	if err != nil {
		details := types.Details{}
		details.VPA = vpa
		val := details.Marshal()
		return check(db.Put(key, val, nil))
	} else {
		var d types.Details
		d.Unmarshal(data)
		d.VPA = vpa
		val := d.Marshal()
		return check(db.Put(key, val, nil))
	}
}

func SetCardDetails(dbName string, key []byte, number, name, expirationMonth, expirationYear, cvv string) error {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	// fetch data from database
	data, err := db.Get(key, nil)
	if err != nil {
		details := types.Details{}
		details.CardNumber = number
		details.CardName = name
		details.CardExpirationMonth = expirationMonth
		details.CardExpirationYear = expirationYear
		details.CardCVV = cvv

		val := details.Marshal()
		return check(db.Put(key, val, nil))
	} else {
		var details types.Details
		details.Unmarshal(data)

		details.CardNumber = number
		details.CardName = name
		details.CardExpirationMonth = expirationMonth
		details.CardExpirationYear = expirationYear
		details.CardCVV = cvv

		val := details.Marshal()
		return check(db.Put(key, val, nil))
	}
}

func GetDetails(dbName string, key []byte) (types.Details, error) {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()
	data, err := db.Get(key, nil)
	if err != nil {
		return types.Details{}, err
	}

	var d types.Details
	d.Unmarshal(data)

	return d, nil

}

func GetRedirectURL(dbName string, key []byte) string {
	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	data, err := db.Get(key, nil)
	if err != nil {
		return ""
	}

	var d types.Details
	d.Unmarshal(data)

	return d.RedirectURL
}

func MoveFromDB(fromDBName string, toDBName string, key []byte) (string, error) {
	fromDB, err := leveldb.OpenFile(fromDBName, nil)
	check(err)
	defer fromDB.Close()

	toDB, err := leveldb.OpenFile(toDBName, nil)
	check(err)
	defer toDB.Close()

	// extract from fromDB
	data, _ := fromDB.Get(key, nil)

	var d types.Details
	d.Unmarshal(data)
	// empty the basket but leave other things intact
	d.Basket = map[string]int{}
	d.Amount = 0

	val := d.Marshal()
	// delete from fromDB
	check(fromDB.Put(key, val, nil))

	// insert into toDB
	// generate unique key
	guid := uuid.New()
	newKey := []byte(guid.String())
	e := check(toDB.Put(newKey, data, nil))
	return guid.String(), e
}

func IterateOverDB(dbName string) {

	db, err := leveldb.OpenFile(dbName, nil)
	check(err)
	defer db.Close()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := string(iter.Key())
		value := string(iter.Value())
		fmt.Println(key, ": ", value)
	}
	iter.Release()
}
