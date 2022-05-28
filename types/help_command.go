package types

import (
	"fmt"
)

// GetAvailableCommands returns all the available commands for telegram channel.
func GetAvailableCommands() string {

	str := "<b><u> List of all Commands </u></b>\n\n"

	str += fmt.Sprintln("<b>/help</b> - <i>Lists all the available commands.</i>\n")
	str += fmt.Sprintln("<b>/initstore</b> - <i>Initialized routine for continuous listening for webapp.</i>\n")
	str += fmt.Sprintln("<b>/save</b> - <i>Alternative to initstore but it saves the data after webapp captures it. And not in realtime.</i>\n")
	str += fmt.Sprintln("<b>/menu</b> - <i>Sends all available items in the menu with a picture, description and price tag.</i>\n")
	str += fmt.Sprintln("<b>/summary</b> - <i>Gives a message with all items in your basket, total price and shipping address, phone number and name.</i>\n")
	str += fmt.Sprintln("<b>/sethouse {housename}</b> - <i>Sets House field of shipping Address to {housename}.</i>\n")
	str += fmt.Sprintln("<b>/setstreet {streetname}</b> - <i>Sets the Street field of Shipping Address to {streetname}.</i>\n")
	str += fmt.Sprintln("<b>/setcity {cityname}</b> - <i>Sets the City field of Shipping Address to {cityname}.</i>\n")
	str += fmt.Sprintln("<b>/setstate {statename}</b> - <i>Sets the State field of Shipping Address to {statename}.</i>\n")
	str += fmt.Sprintln("<b>/setcountry {countryname}</b> - <i>Sets the Country field of Shipping Address to {countryname}.</i>\n")
	str += fmt.Sprintln("<b>/setpostcode {postcode}</b> - <i>Sets the Postcode field of Shipping Address to {postcode}.</i>\n")
	str += fmt.Sprintln("<b>/setname {name}</b> - <i>Sets the name field of Details to {name}.</i>\n")
	str += fmt.Sprintln("<b>/setcard {cardNumber}, {cardName}, {expirationMonth}, {expirationYear}, {cvv}</b> - <i>Sets the Card details in Card field of Details.</i>\n")
	str += fmt.Sprintln("<b>/card {cardNumber}, {cardName}, {expirationMonth}, {expirationYear}, {cvv}</b> - <i>This is same as setcard but when issued it also initiates checkout.</i>\n")
	str += fmt.Sprintln("<b>/setvpa {vpa}</b> - <i>Sets the VPA field of Details to {vpa}</i>\n")
	str += fmt.Sprintln("<b>/vpa {vpa}</b> - <i>This is same as setvpa but when issued it also starts checkout.</i>\n")
	str += fmt.Sprintln("<b>/checkout</b> - <i>Proceeds for checkout.</i>\n")
	return str
}
