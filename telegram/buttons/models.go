package buttons

type KeyboardButton struct {
	Type string `json:"-"`
	Text string `json:"text"`
}

func BasicButton(text string) KeyboardButton{
	return KeyboardButton{
		Type: "keyboardButton",
		Text: text,
	}
}
func RequestGeoButton(text string) KeyboardButton{
	return KeyboardButton{
		Type: "keyboardButtonRequestGeoLocation",
		Text: text,
	}
}