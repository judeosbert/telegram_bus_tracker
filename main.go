package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/judeosbert/bus_tracker_bot/models"
)

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	_, err := parseTelegramMessageRequest(r)
	if err != nil {
		log.Printf("Failed to parse telegram request,%s", err.Error())
		return
	}
}

func parseTelegramMessageRequest(r *http.Request) (*models.Update, error) {
	var update models.Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return nil, err
	}
	return &update, nil
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
