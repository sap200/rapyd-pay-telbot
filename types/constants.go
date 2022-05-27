package types

const Telegram_Bot_Token = "<<Telegram-Bot-Token>>"
const Base_url = "https://sandboxapi.rapyd.net"
const Access_key = "<<Rapyd-Sandbox-Access-Key>>"
const Secret_key = "<<Rapyd-Sandbox-Secret-Key>>"

// db containing basket path
const DB_Details_Path = "./storage/detailsDB"
const DB_Orders_Path = "./storage/ordersDB"

// add or remove to basket command
const ADD_TO_BASKET = "add-"
const REMOVE_FROM_BASKET = "remove-"
const PAYMENT_CONFIRM = "payconfirm-"
const PAYMENT_REJECT = "payreject-"
const GIBBERISH = "gibberish"

// shipping address fields
const HOUSE_FIELD = "house"
const STREET_FIELD = "street"
const CITY_FIELD = "city"
const STATE_FIELD = "state"
const COUNTRY_FIELD = "country"
const POSTCODE_FIELD = "postcode"
const PHONENUMBER_FIELD = "phone"
const NAME_FIELD = "name"

// rapyd payment questions
const PAYMENT_METHODS = "choose a payment method"
const PAYMENT_CATEGORY = "choose a payment category"
const PAYMENT_DELIMITER = "-"
const PAYMENT_ACTIVE = "ACT"
const PAYMENT_CLOSED = "CLO"
const PAYMENT_ERROR = "ERR"

// Payment Types
const BANK = "Bank"
const UPI = "UPI"
const CARD = "Card"
const EWALLET = "eWallet"
