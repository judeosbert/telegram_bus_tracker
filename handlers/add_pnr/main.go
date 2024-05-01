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
		requestForPnr(stateSaver, bot, chatId)
	default:
		handlers.SendMessage(bot, chatId, "Error. Start over with commands")
	}

}

func requestForPnr(stateSaver state.Saver, bot *telego.Bot, chatId int64) {
	stateSaver.SetUserState(chatId, Init{})
	handlers.SendMessage(bot, chatId, "Okay, Send Trip Code")
}

var Predicate = th.CommandEqual("add_pnr")
