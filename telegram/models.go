package telegram

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
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
	IsPresistent bool `json:"is_persistent"`
	ResizeKeyboard bool `json:"resize_keyboard"` 
	OneTimeKeyboard bool `json:"one_time_keyboard"`
	InputFieldPlaceholder string `json:"input_field_placeholder"`
	Selective bool `json:"selective"`
}

type KeyboardButton struct {
	Text string `json:"text"`
	RequestLocation bool `json:"request_location"`
}

