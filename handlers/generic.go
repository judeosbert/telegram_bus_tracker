package handlers

import (
	"log"

	"github.com/mymmrac/telego"
)

func SendMessage(bot *telego.Bot, chatId int64,msg string) error {
	sentMessage, err := bot.SendMessage(
		&telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "Okay,Send PNR",
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sent message %+v\n", sentMessage)
	return nil
}
