package botengine

import (
	"errors"

	"github.com/judeosbert/bus_tracker_bot/admin"
	addpnr "github.com/judeosbert/bus_tracker_bot/handlers/add_pnr"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
)

func HandleAfterPnrSent(message telego.Message, saver state.Saver,admin admin.AdminUtils) (*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetState(string(chatId))
	if err != nil {
		return nil, errors.New("Wrong State")
	}
	switch prevState.(type) {
	case addpnr.Init:
		pnr := message.Text
		saver.SetState(string(chatId), addpnr.SavePnr{
			Pnr: pnr,
		})
		return &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text:            "Which bus service have you booked?",
			ReplyParameters: &telego.ReplyParameters{},
			ReplyMarkup: &telego.ReplyKeyboardMarkup{
				Keyboard: [][]telego.KeyboardButton{
					{
						telego.KeyboardButton{
							Text: "Kerala SRTC",
						},
						telego.KeyboardButton{
							Text: "Karnataka SRTC",
						},
					},
				},
				IsPersistent:          false,
				ResizeKeyboard:        false,
				OneTimeKeyboard:       true,
				InputFieldPlaceholder: "Select your bus service provider",
				Selective:             false,
			},
		}, nil

	default:
		return nil, errors.New("Wrong state")
	}
}

func HandleAfterProviderSent(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) (*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetState(string(chatId))
	if err != nil {
		return nil, errors.New("Wrong State")
	}
	switch prevState:= prevState.(type) {
	case addpnr.SavePnr:
		provider := message.Text
		if(provider == ""){
			return nil,errors.New("Empty Provider")
		}
		saver.SetState(string(chatId), addpnr.SavePnrProvider{
			Pnr:             prevState.Pnr,
			ServiceProvider: provider,
		})
		adminUtils.SendForVerification(admin.NewTripInfo{
			ServiceProvider: provider,
			Pnr:             prevState.Pnr,
		})
		return &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "Okay, PNR wll be verifed and added. Check after some time. ",
		}, nil

	default:
		return nil, errors.New("Wrong state")
	}
}
