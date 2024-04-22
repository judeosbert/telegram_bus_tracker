package main

import (
	"log"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
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
	botToken := "7073126054:AAEI729OK0391qRMrXzpojWqB-5ROuwPi_I"
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}

	// err = bot.SetWebhook(&telego.SetWebhookParams{
	// 	URL: "https://telegrambustracker-production.up.railway.app/bot" + botToken,
	// })
	// if err != nil {
	// 	log.Printf("Failed to set webhook info %s", err.Error())
	// 	return
	// }

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Printf("Failed to get webhook info %s", err.Error())
		return
	}
	log.Printf("Webhook Info %+v/n", info)

	updates, err:= bot.UpdatesViaWebhook("/bot" + bot.Token())
	if(err != nil ){
		log.Printf("Failed to get updates via hook %s",err.Error())
		return
	}

	go func() {
		bot.StartWebhook(":8080")
	}()
	defer func() {
		bot.StopWebhook()
	}()
	for update := range updates {
		if update.Message == nil {
			log.Println("Update is Empty")
			continue
		}
		chatId := update.Message.Chat.ID
		sentMessage, err := bot.SendMessage(telegoutil.Message(
			telegoutil.ID(chatId),
			"Hello World",
		))
		if err != nil {
			log.Printf("Failed to send message %s", err.Error())
			continue
		}
		log.Printf("Sent Message: %v\n", sentMessage)
	}

}
