package botengine

import (
	"errors"

	"github.com/judeosbert/bus_tracker_bot/admin"
	addpnr "github.com/judeosbert/bus_tracker_bot/handlers/add_pnr"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
)

func HandleAfterPnrSend(message telego.Message, saver state.Saver, admin admin.AdminUtils) (*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(string(chatId))
	if err != nil {
		return nil, errors.New("Wrong State")
	}
	switch prevState.(type) {
	case addpnr.Init:
		pnr := message.Text
		saver.SetUserState(string(chatId), addpnr.SavePnr{
			Pnr: pnr,
		})
		return &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "What is your trip code?",
		}, nil

	default:
		return nil, errors.New("Wrong state")
	}
}

func HandleAfterPnrTripCodeSent(message telego.Message, saver state.Saver, admin admin.AdminUtils) (*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(string(chatId))
	if err != nil {
		return nil, errors.New("Wrong State")
	}
	switch prevState := prevState.(type) {
	case addpnr.SavePnr:
		tripCode := message.Text
		tripState, err := saver.GetTripState(tripCode)
		if err == nil {
			switch tripState.(type) {
			case addpnr.SubmittedForVerification:
				return &telego.SendMessageParams{
					ChatID:               telego.ChatID{
						ID:chatId,
					},
					Text:                 "The trip is submitted for verification. You will be notified. Or your can check later",
				},nil
				saver.AddTripObserver(tripCode,string(chatId))
			}
			return nil, errors.New("Existing state for tripcode "+tripCode) 
		}
		saver.SetUserState(string(chatId), addpnr.SavePnrTripCode{
			Pnr:      prevState.Pnr,
			TripCode: tripCode,
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
	prevState, err := saver.GetUserState(string(chatId))
	if err != nil {
		return nil, errors.New("Wrong State")
	}
	switch prevState := prevState.(type) {
	case addpnr.SavePnrTripCode:
		provider := message.Text
		saver.SetUserState(string(chatId), addpnr.SavePnrTripCodeProvider{
			Pnr:             prevState.Pnr,
			TripCode:        prevState.TripCode,
			ServiceProvider: provider,
		})
		adminUtils.SendForVerification(admin.NewTripInfo{
			ServiceProvider: provider,
			TripCode:        prevState.TripCode,
			Pnr:             prevState.Pnr,
		})

		saver.SetTripState(prevState.TripCode,addpnr.SubmittedForVerification{
			Pnr:             prevState.Pnr,
			TripCode:        prevState.TripCode,
			ServiceProvider: provider,
		})

		saver.AddTripObserver(prevState.TripCode,string(chatId))
		
		return &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "Okay, PNR wll be verifed and added. Check after some time. You will be notified",
		}, nil

	default:
		return nil, errors.New("Wrong state")
	}
}
