package main

import (
	"fmt"
	"log"
	"os"
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
	stateSaver := state.NewStateSaver()
	botToken := os.Getenv("TELEGO_BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}
	bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands: []telego.BotCommand{
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
	admin := admin.NewAdminUtils(bot, stateSaver)
	botEngine := botengine.NewBotEnginer(stateSaver, admin)

	var env = os.Getenv("MODE")
	var updates <-chan telego.Update
	if env == "dev" {
		updates, err = bot.UpdatesViaLongPolling(
			&telego.GetUpdatesParams{
				Offset:  0, // Will be automatically updated by UpdatesViaLongPolling
				Timeout: 8, // Can be set instead of using WithLongPollingUpdateInterval (default, recommended way)
			}, telego.WithLongPollingUpdateInterval(time.Second*0), telego.WithLongPollingRetryTimeout(time.Second*8), telego.WithLongPollingBuffer(100))
		if err != nil {
			panic(err)
		}

		defer bot.StopLongPolling()
	} else {
		_ = bot.SetWebhook(&telego.SetWebhookParams{
			URL: "https://telegrambustracker-production.up.railway.app/bot",
		})

		info, _ := bot.GetWebhookInfo()
		fmt.Printf("Webhook Info: %+v\n", info)

		updates, _ = bot.UpdatesViaWebhook("bot")

		go func() {
			_ = bot.StartWebhook("localhost:443")
		}()

		// Stop reviving updates from update channel and shutdown webhook server
		defer func() {
			_ = bot.StopWebhook()
		}()
	}

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		panic(err)
	}
	defer bh.Stop()

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		addpnr.Handler(bot, update, stateSaver)
	}, addpnr.Predicate)

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		start.Handler(bot, update, stateSaver)
	}, start.Predicate)

	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		deletetrip.Handler(bot, update, stateSaver)
	}, deletetrip.Predicate)

	bh.HandleMessage(func(bot *telego.Bot, message telego.Message) {
		botEngine.PostUpdate(message)
	})

	bh.HandleCallbackQuery(func(bot *telego.Bot, query telego.CallbackQuery) {

	}, th.AnyCallbackQuery())

	go func() {
		for msg := range botEngine.OutChan() {
			log.Printf("Sending message out from channel %+v\n", msg)
			bot.SendMessage(&msg)
		}
	}()
	go func() {
		botEngine.Start()
	}()
	bh.Start()
}
