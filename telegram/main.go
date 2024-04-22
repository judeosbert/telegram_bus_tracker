package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/judeosbert/bus_tracker_bot/telegram/buttons"
)

func HandleIncomingMessage(context *gin.Context) {
	// var update map[string]map[string]interface{}
	// err := context.BindJSON(&update)
	// if err != nil {
	// 	fmt.Printf("error %s", err.Error())
	// 	return
	// }
	// fmt.Println("Body %s", update["message"]["text"])
	update := &Update{}
	err := context.Bind(update)
	if err != nil {
		replyWithText(update.Message.Chat.ID, err.Error())
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	msg, err := handleMessage(*update)
	if err != nil {
		replyWithText(update.Message.Chat.ID, err.Error())
		return
	}

	sendMessageToTelegramChat(*msg)

}

func handleMessage(update Update) (*ReplyMessage, error) {
	keyboard := make([]buttons.KeyboardButton,0)
	keyboard = append(keyboard, buttons.BasicButton("Sample Button"))
	keyboard = append(keyboard, buttons.RequestGeoButton("Send Location"))

	return &ReplyMessage{
		ChatId:  strconv.Itoa(update.Message.Chat.ID),
		Message: "Sample Message",
		ReplyMarkup: ReplyKeyboardMarkup{
			Keyboard:              keyboard,
			OneTimeKeyboard:       true,
			InputFieldPlaceholder: "",
			Selective:             false,
		},
	}, nil
}

func replyWithKeyboard(chatId int, keyboardOptions ReplyKeyboardMarkup, message string) {
	sendMessageToTelegramChat(ReplyMessage{
		ChatId:      strconv.Itoa(chatId),
		Message:     message,
		ReplyMarkup: keyboardOptions,
	})
}

func replyWithText(chatId int, text string) (string, error) {
	return sendMessageToTelegramChat(ReplyMessage{
		ChatId:  strconv.Itoa(chatId),
		Message: text,
	})
}

func sendMessageToTelegramChat(reply ReplyMessage) (string, error) {
	telegramEp := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", "7073126054:AAEI729OK0391qRMrXzpojWqB-5ROuwPi_I")
	// response, err := http.PostForm(
	// 	telegramEp,

	// )

	// if err != nil {
	// 	log.Printf("Failed to post message to chat %s", err.Error())
	// 	return "", err
	// }

	// defer response.Body.Close()

	// bodyBytes, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	log.Printf("error parsing telegram %s", err.Error())
	// 	return "", err
	// }
	// bodyString := string(bodyBytes)
	// log.Printf("Body of telegram message post %s", bodyString)
	// return bodyString, nil
	buf, err := json.Marshal(reply)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", telegramEp, bytes.NewBuffer(buf))
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resBuf := []byte{}
	_,err = res.Body.Read(resBuf)
	if err != nil {
		return "",err
	}
	

	if res.StatusCode == 200 {
		log.Printf("Response Body %s",resBuf)
		return string(resBuf) , nil
	}

	return "error", errors.New("failed to send request:Unknown reason")
}
