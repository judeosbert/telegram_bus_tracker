package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleIncomingMessage(context *gin.Context) {
	var update Update
	err := context.BindJSON(&update)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	msg, err := handleMessage(update)
	if err != nil {
		sendMessageToTelegramChat(update.Message.Chat.Id, err.Error())
		return
	}

	sendMessageToTelegramChat(update.Message.Chat.Id, msg.Text)

}

func handleMessage(update Update) (*Message, error) {
	return &Message{
		Text: "Sample",
		Chat: update.Message.Chat,
	}, nil
}

func sendMessageToTelegramChat(chatId int, text string) (string, error) {
	telegramEp := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", "7073126054:AAEI729OK0391qRMrXzpojWqB-5ROuwPi_I")
	response, err := http.PostForm(
		telegramEp,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			text:      {text},
		})

	if err != nil {
		log.Printf("Failed to post message to chat %s", err.Error())
		return "", err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error parsing telegram %s", err.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of telegram message post %s", bodyString)
	return bodyString, nil
}