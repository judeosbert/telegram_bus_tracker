package telegram

import "github.com/judeosbert/bus_tracker_bot/telegram/buttons"

type Update struct {
	Message Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int    `json:"date"`
	Text      string `json:"text"`
}

type Chat struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type From struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type ReplyMessage struct {
	Flags int `json:"flags omitempty"`
	NoWebPage bool `json:"no_webpage omitempty"`
	Silent bool `json:"silent omitempty"`
	Background bool `json:"background omitempty"`
	ClearDraft bool `json:"clear_draft omitempty"`
	NoForwards bool `json:"noforwards omitempty"`
	InvertMedia bool `json:"invert_media omitempty"`
	ChatId string `json:"chat_id"`
	Message string `json:"message omitempty"`
	ReplyMarkup ReplyKeyboardMarkup `json:"replyKeyboardMarkup omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard []buttons.KeyboardButton `json:"rows"`
	IsPresistent bool `json:"persistent"`
	ResizeKeyboard bool `json:"resize"` 
	OneTimeKeyboard bool `json:"single_use"`
	InputFieldPlaceholder string `json:"placeholder"`
	Selective bool `json:"selective"`
}






