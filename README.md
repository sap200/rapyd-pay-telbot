# rapyd-pay-telbot

![Project Image](alex-kinght.jpg "Rapyd Pay Telbot")


### Rapyd Pay Telbot is a telegram frontend checkout experience for users in effort to build galaxy's best checkout experience.

## Usage

[Try it out](https://telegram.me/RapydPBot)


## List of available commands

**/help** - Lists all the available commands.

**/initstore** - Initialized routine for continuous listening for webapp.

**/initstore** - Alternative to initstore but it saves the data after webapp captures it. And not in realtime.

**/menu** - Sends all available items in the menu with a picture, description and price tag.

**/summary** - Gives a message with all items in your basket, total price and shipping address, phone number and name.

**/sethouse {housename}** - Sets House field of shipping Address to {housename}.

**/setstreet {streetname}** - Sets the Street field of Shipping Address to {streetname}.

**/setcity {cityname}** - Sets the City field of Shipping Address to {cityname}.

**/setstate {statename}** - Sets the State field of Shipping Address to {statename}.

**/setcountry {countryname}** - Sets the Country field of Shipping Address to {countryname}.

**/setpostcode {postcode}** - Sets the Postcode field of Shipping Address to {postcode}.

**/setname {name}** - Sets the name field of Details to {name}.

**/setcard {cardNumber}, {cardName}, {expirationMonth}, {expirationYear}, {cvv}** - Sets the Card details in Card field of Details.

**/card {cardNumber}, {cardName}, {expirationMonth}, {expirationYear}, {cvv}** - This is same as setcard but when issued it also initiates checkout.

**/setvpa {vpa}** - Sets the VPA field of Details to {vpa}

**/vpa {vpa}** - This is same as setvpa but when issued it also starts checkout.

**/checkout** - Proceeds for checkout
