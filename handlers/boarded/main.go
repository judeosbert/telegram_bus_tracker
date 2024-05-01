package boarded

import (

	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Handler(bot *telego.Bot, update telego.Update, stateSaver state.Saver) {
	chatId := update.Message.Chat.ID
	_, err := stateSaver.GetActiveTrip(chatId)
	if err != nil {
		handlers.SendMessage(bot, chatId, "You have no active trip.")
		return
	}
	// if err :=stateSaver.MarkBoarded(chatId);err != nil {
	// 	handlers.SendMessage(bot,chatId,"Error marking you as boarded. Try again later.")
	// 	return
	// }

	bot.SendMessage(&telego.SendMessageParams{
		ChatID: telego.ChatID{
			ID: chatId,
		},
		Text:               "Wonderful! Please share your current location with me.",
		LinkPreviewOptions: &telego.LinkPreviewOptions{},
		ReplyMarkup: &telego.ReplyKeyboardMarkup{
			Keyboard: [][]telego.KeyboardButton{
				{
					{
				
						Text:            "Send Current Location",
						RequestLocation: true,
					},
				},
			},
			OneTimeKeyboard:       true,
			InputFieldPlaceholder: "Send your current location",
			Selective:             false,
		},
	})
}

var Predicate = th.CommandEqual("boarded")
