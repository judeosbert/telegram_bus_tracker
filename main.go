package main

import (
	"log"
	"time"

	"github.com/judeosbert/bus_tracker_bot/admin"
	botengine "github.com/judeosbert/bus_tracker_bot/bot_engine"
	addpnr "github.com/judeosbert/bus_tracker_bot/handlers/add_pnr"
	deletetrip "github.com/judeosbert/bus_tracker_bot/handlers/delete_trip"
	"github.com/judeosbert/bus_tracker_bot/handlers/start"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func main() {
	// 	r := gin.Default()
	// 	r.POST("/", func(ctx *gin.Context) {
	// 		body, _ := io.ReadAll(ctx.Request.Body)
	// 		fmt.Printf("incoming body   %s",string(body))

	// 		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
	// 		telegram.HandleIncomingMessage(ctx)
	// 	})
	// 	r.Run()

	stateSaver := state.NewStateSaver()
	
	botToken := "7073126054:AAEI729OK0391qRMrXzpojWqB-5ROuwPi_I"
	
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}
	bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands:     []telego.BotCommand{
			{
				Command:     "/start",
				Description: "Copy paste ticket",
			},
			{
				Command:     "/add_trip_manual",
				Description: "Manually Add a new Trip",
			},
		},
	})
	admin := admin.NewAdminUtils(bot,stateSaver)
	botEngine := botengine.NewBotEnginer(stateSaver,admin)

	updates, err := bot.UpdatesViaLongPolling(
		&telego.GetUpdatesParams{
			Offset:  0, // Will be automatically updated by UpdatesViaLongPolling
			Timeout: 8, // Can be set instead of using WithLongPollingUpdateInterval (default, recommended way)
		}, telego.WithLongPollingUpdateInterval(time.Second*0), telego.WithLongPollingRetryTimeout(time.Second*8), telego.WithLongPollingBuffer(100))
	if err != nil {
		panic(err)
	}

	defer bot.StopLongPolling()

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		panic(err)
	}
	defer bh.Stop()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		addpnr.Handler(bot,update,stateSaver)
	},addpnr.Predicate)
	
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		start.Handler(bot,update,stateSaver)
	},start.Predicate)
	
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		deletetrip.Handler(bot,update,stateSaver)
	},deletetrip.Predicate)

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		botEngine.PostUpdate(message)
	})
	
	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {
		
	},th.AnyCallbackQuery())

	
	go func(){
		for msg := range botEngine.OutChan(){
			log.Printf("Sending message out from channel %+v\n",msg)
			bot.SendMessage(&msg)
		}
	}()
	go func(){
		botEngine.Start()
	}()
	bh.Start()
}
