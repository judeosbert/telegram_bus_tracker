package addpnr

import (
	"log"

	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Handler(bot *telego.Bot, update telego.Update, stateSaver state.Saver) {
	message := update.Message
	if message == nil {
		log.Println("Empty Message")
		return
	}

	chatId := message.Chat.ID
	prevState, _ := stateSaver.GetUserState(chatId)

	switch prevState {
	case nil:
		requestForBusProvider(stateSaver, bot, chatId)
	default:
		stateSaver.SetUserState(chatId, nil)
		handlers.SendMessage(bot, chatId, "Error. Start over with commands")
	}

}

func requestForBusProvider(stateSaver state.Saver, bot *telego.Bot, chatId int64) {

	stateSaver.SetUserState(chatId, Init{})
	msg := &telego.SendMessageParams{
		ChatID:               telego.ChatID{
			ID:       chatId,
		},
		Text:                 "Okay, Select your bus service provider",
		ReplyMarkup: &telego.ReplyKeyboardMarkup{
			Keyboard:              [][]telego.KeyboardButton{
				{
					{
						Text: "Kerala SRTC",
					},
					{
						Text: "Karnataka SRTC",
					},
				},
			},
			IsPersistent:          false,
			ResizeKeyboard:        false,
			OneTimeKeyboard:       true,
			InputFieldPlaceholder: "Select your bus service provider",
			Selective:             false,
		},
	}
	bot.SendMessage(msg)
}

var Predicate = th.CommandEqual("add_trip_manual")
