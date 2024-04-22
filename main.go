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
	bot,err := telego.NewBot(botToken,telego.WithDefaultDebugLogger())
	if err != nil {
		panic(err)
	}
	updates, _ := bot.UpdatesViaWebhook("/")

	go func ()  {
		bot.StartWebhook("0.0.0.0:443")
	}()
	defer func(){
		bot.StopWebhook()
	}()
	for update := range updates{
		if(update.Message == nil){
			log.Println("Update is Empty")
			continue
		}
		chatId := update.Message.Chat.ID
		sentMessage, err := bot.SendMessage(telegoutil.Message(
			telegoutil.ID(chatId),
			"Hello World",

		))
		if(err != nil){
			log.Printf("Failed to send message %s",err.Error())
			continue
		}
		log.Printf("Sent Message: %v\n", sentMessage)
	}


}
