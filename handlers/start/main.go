package start

import (
	"log"

	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
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
		requestForTicket(stateSaver, bot, chatId)
	default:
		stateSaver.SetUserState(chatId, nil)
		handlers.SendMessage(bot, chatId, "Error. Start over with commands")
	}

}

func requestForTicket(stateSaver state.Saver, bot *telego.Bot, chatId int64) {
	stateSaver.SetUserState(chatId, RequestTicket{})
	bot.SendMessage(&telego.SendMessageParams{
		ChatID:               tu.ID(chatId),
		Text:                 "Copy paste the ticket here with bus no.",
	})
}

var Predicate = th.CommandEqual("start")
