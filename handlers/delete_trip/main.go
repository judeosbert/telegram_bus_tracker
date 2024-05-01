package deletetrip

import (
	"fmt"
	"log"

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
	if err :=stateSaver.DeleteActiveTrip(chatId);err != nil {
		handlers.SendMessage(bot,chatId,"Error removing you from this trip. Try again later.")
		return
	}

	//remove from group.
	err = bot.BanChatMember(&telego.BanChatMemberParams{
		ChatID:         telego.ChatID{ID: chatId},
		UserID:         chatId,
		RevokeMessages: true,
	})
	if(err != nil){
		log.Println(fmt.Sprintf("Could not kick %s from %s. Need to retry",chatId,trip))
	}
	handlers.SendMessage(bot,chatId,fmt.Sprintf("You are removed from %s",trip))
}

var Predicate = th.CommandEqual("delete_trip")
