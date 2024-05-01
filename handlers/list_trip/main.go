package listtrip

import (
	"fmt"

	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func Handler(bot *telego.Bot, update telego.Update, stateSaver state.Saver) {
	chatId := update.Message.Chat.ID
	trip,err := stateSaver.GetActiveTrip(chatId)
	if err != nil {
		handlers.SendMessage(bot,chatId,"You have no active trip.")
		return
	}
	handlers.SendMessage(bot,chatId,fmt.Sprintf("You are part of the trip %s. Use /add_pnr to get group link(if available)",trip))
}

var Predicate = th.CommandEqual("delete_trip")
