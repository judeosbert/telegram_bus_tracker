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
	prevState, _ := stateSaver.GetUserState(string(chatId))

	switch prevState {
	case nil:
		requestForPnr(stateSaver, bot, chatId)
	default:
		handlers.SendMessage(bot, chatId, "Error. Start over with commands")
	}

}

func requestForPnr(stateSaver state.Saver, bot *telego.Bot, chatId int64) {
	stateSaver.SetUserState(string(chatId), Init{})
	handlers.SendMessage(bot, chatId, "Send the PNR")
}

var Predicate = th.CommandEqual("add_pnr")
