package admin

import (
	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
)

func HandleCallbackQuery(bot *telego.Bot, query telego.CallbackQuery, admin AdminUtils, saver state.Saver) {
	chatId := query.From.ID
	// err := HandleTripStateVerficationCallback(bot, chatId, query.Data, admin, saver)
	// if err == nil {
	// 	return
	// }
	//Add other callback handlers
	handlers.SendMessage(bot, chatId, "Unknown operation")

}
