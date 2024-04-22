package buttons

type KeyboardButton struct {
	Text string `json:"text"`
	RequestLocation bool `json:"request_location,omitempty"`

}

func BasicButton(text string) KeyboardButton{
	return KeyboardButton{
		Text: text,
	}
}
func RequestGeoButton(text string) KeyboardButton{
	return KeyboardButton{
		Text: text,
		RequestLocation: true,
	}
}