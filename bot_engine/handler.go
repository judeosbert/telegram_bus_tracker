package botengine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/judeosbert/bus_tracker_bot/admin"
	addpnr "github.com/judeosbert/bus_tracker_bot/handlers/add_pnr"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/mymmrac/telego"
)

func HandleAfterTripCodeSend(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	tripCode := message.Text
	prevTripState, _ := saver.GetTripState(tripCode)
	if prevTripState != nil {
		switch prevTripState := prevTripState.(type) {
		case addpnr.SubmittedForVerification:
			trip, err := saver.GetActiveTrip(chatId)
			if err == nil {
				return []*telego.SendMessageParams{
					{
						ChatID: telego.ChatID{
							ID: chatId,
						},
						Text: fmt.Sprintf("You are already part of %s. Either /delete_trip to delete the active trip or wait until its complete.", trip),
					},
				}, nil
			}
			saver.SetActiveTrip(chatId, tripCode)
			saver.AddTripObserver(tripCode, chatId)
			return []*telego.SendMessageParams{{
				ChatID: telego.ChatID{
					ID: chatId,
				},
				Text: "The trip is submitted for verification. You will be notified",
			}}, nil
		case admin.TripStateValidation:
			if prevTripState.Status == admin.STATUS_REJECTED_TRIP_VERIFICATION {
				return []*telego.SendMessageParams{{
					ChatID: telego.ChatID{
						ID: chatId,
					},
					Text: "The tripcode was rejected for incorrect data. ",
				}}, nil
			}

			if prevTripState.Status == admin.STATUS_VERIFIED_TRIP_VERIFICATION {
				trip, err := saver.GetActiveTrip(chatId)
				if err == nil {
					return []*telego.SendMessageParams{
						{
							ChatID: telego.ChatID{
								ID: chatId,
							},
							Text: fmt.Sprintf("You are already part of %s. Either /delete_trip to delete the active trip or wait until its complete.", trip),
						},
					}, nil
				}
				saver.SetActiveTrip(chatId, tripCode)
				saver.AddTripObserver(tripCode, chatId)
				return []*telego.SendMessageParams{{
					ChatID: telego.ChatID{
						ID: chatId,
					},
					Text: "The trip is verified. The trip group will be shared soon. You will be notified",
				}}, nil
			}

		case admin.TripStateVerifiedWithLink:
			return []*telego.SendMessageParams{{
				ChatID: telego.ChatID{
					ID: chatId,
				},
				Text: fmt.Sprintf("The trip is verified. Join this group for update %s", prevTripState.InviteLink),
			}}, nil

		}
	}
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("wrong state")
	}
	switch prevState.(type) {
	case addpnr.Init:
		saver.SetUserState(chatId, addpnr.SaveTripCode{
			TripCode: tripCode,
		})
		return []*telego.SendMessageParams{{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "What is your pnr?",
		}}, nil

	default:
		return nil, errors.New("wrong state")
	}
}

func HandleAfterTripCodePnrSent(message telego.Message, saver state.Saver, admin admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("wrong state")
	}
	switch prevState := prevState.(type) {
	case addpnr.SaveTripCode:
		pnr := message.Text
		saver.SetUserState(chatId, addpnr.SaveTripCodePnr{
			Pnr:      pnr,
			TripCode: prevState.TripCode,
		})
		return []*telego.SendMessageParams{
			{
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
			},
		}, nil

	default:
		return nil, errors.New("wrong state")
	}
}

func HandleAfterProviderSent(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("wrong state")
	}
	switch prevState := prevState.(type) {
	case addpnr.SaveTripCodePnr:
		provider := message.Text
		saver.SetUserState(chatId, addpnr.SaveTripCodePnrProvider{
			Pnr:             prevState.Pnr,
			TripCode:        prevState.TripCode,
			ServiceProvider: provider,
		})
		adminUtils.SendForVerification(admin.NewTripInfo{
			ServiceProvider: provider,
			TripCode:        prevState.TripCode,
			Pnr:             prevState.Pnr,
		})

		saver.SetTripState(prevState.TripCode, addpnr.SubmittedForVerification{
			Pnr:             prevState.Pnr,
			TripCode:        prevState.TripCode,
			ServiceProvider: provider,
		})
		trip, err := saver.GetActiveTrip(chatId)
		if err == nil {
			return []*telego.SendMessageParams{
				{
					ChatID: telego.ChatID{
						ID: chatId,
					},
					Text: fmt.Sprintf("You are already part of %s. Either /delete_trip to delete the active trip or wait until its complete.", trip),
				},
			}, nil
		}
		saver.SetActiveTrip(chatId, prevState.TripCode)
		saver.AddTripObserver(prevState.TripCode, chatId)

		return []*telego.SendMessageParams{{
			ChatID: telego.ChatID{
				ID: chatId,
			},
			Text: "Okay, your trip and pnr will be verifed and added. You will be notified.",
		}}, nil

	default:
		return nil, errors.New("wrong state")
	}
}

func HandleInviteLinkMsg(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	if message.ReplyToMessage == nil {
		return nil, errors.New("Not a reply to a message")
	}
	orgMsg := message.ReplyToMessage.Text
	if len(orgMsg) == 0 {
		return nil, errors.New("Original Message is empty")
	}
	if !strings.Contains(orgMsg, "#VERIFICATION") {
		return nil, errors.New("Not a reply to verification message")
	}
	inviteLink := message.Text
	if len(inviteLink) == 0 {
		return nil, errors.New("Empty Invite Link")
	}
	r, _ := regexp.Compile("(TripCode:)[A-Z0-9]+")
	p := strings.Split(r.FindString(orgMsg), ":")
	if len(p) != 2 {
		return nil, errors.New(fmt.Sprintf("Incorrect Match %+v", p))
	}
	tripCode := p[1]

	r, _ = regexp.Compile("(PNR:)[A-Z0-9]+")
	p = strings.Split(r.FindString(orgMsg), ":")
	if len(p) != 2 {
		return nil, errors.New(fmt.Sprintf("Incorrect Match %+v", p))
	}
	pnr := p[1]

	saver.SetTripState(tripCode, &admin.TripStateVerifiedWithLink{
		Pnr:        pnr,
		TripCode:   tripCode,
		InviteLink: inviteLink,
	})

	obs := saver.GetTripObservers(tripCode)
	if len(obs) == 0 {
		return nil, errors.New("No observers")
	}
	msgs := []*telego.SendMessageParams{}
	for i := 0; i < len(obs); i++ {
		ob := obs[i]

		msgs = append(msgs, &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: int64(ob),
			},
			Text: fmt.Sprintf("Join this group for trip updates\n%s", inviteLink),
			LinkPreviewOptions: &telego.LinkPreviewOptions{
				IsDisabled:       false,
				URL:              inviteLink,
				PreferLargeMedia: true,
				ShowAboveText:    true,
			},
		})
	}
	return msgs, nil
}

func HandleTripStateVerfication(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	if message.ReplyToMessage == nil {
		return nil, errors.New("Not a reply to a message")
	}
	orgMsg := message.ReplyToMessage.Text
	if len(orgMsg) == 0 {
		return nil, errors.New("Original Message is empty")
	}
	if !strings.Contains(orgMsg, "#VERIFICATION") {
		return nil, errors.New("Not a reply to verification message")
	}

	r, _ := regexp.Compile("(TripCode:)[A-Z0-9]+")
	p := strings.Split(r.FindString(orgMsg), ":")
	if len(p) != 2 {
		return nil, errors.New(fmt.Sprintf("Incorrect Match %+v", p))
	}
	tripCode := p[1]

	r, _ = regexp.Compile("(PNR:)[A-Z0-9]+")
	p = strings.Split(r.FindString(orgMsg), ":")
	if len(p) != 2 {
		return nil, errors.New(fmt.Sprintf("Incorrect Match for PNR %+v", p))
	}
	pnr := p[1]

	state := message.Text

	var update string
	if state == "#Verified" {
		saver.SetTripState(tripCode, admin.NewStateTripVerified(pnr, tripCode))
		update = fmt.Sprintf("Your trip %s is verified. Join the channel for updates.", tripCode)

	} else if state == "#Rejected" {
		saver.SetTripState(tripCode, admin.NewStateTripRejected(pnr, tripCode))
		update = fmt.Sprintf("Your trip %s is rejected. Join the channel for updates.", tripCode)
	} else {
		return nil, errors.New("Unknown verification status msg")
	}

	obs := saver.GetTripObservers(tripCode)
	msgs := []*telego.SendMessageParams{}
	for i := 0; i < len(obs); i++ {
		chatId := obs[i]
		msgs = append(msgs, &telego.SendMessageParams{
			ChatID: telego.ChatID{
				ID: int64(chatId),
			},
			Text: update,
		})

	}

	return msgs, nil
}

func HandleGroupChatMsgCommands(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	if message.Chat.Type != "supergroup" {
		return nil, errors.New("chat is not from group")
	}
	if message.Entities == nil {
		return nil, errors.New("No entities")
	}
	var hasCommand = false
	for i := 0; i < len(message.Entities); i++ {
		e := message.Entities[i]
		if e.Type == "bot_command" {
			hasCommand = true
			break
		}
	}
	if !hasCommand {
		return nil, errors.New("not a command message")
	}
	tripCode := message.Chat.Title
	if len(tripCode) == 0 {
		return nil, errors.New("trip code not found in group title")
	}
	chatId := message.Chat.ID
	saver.SetTripGroup(tripCode, chatId)

	//supported commands get update.
	ps := strings.Split(message.Text, "@")
	command := ps[0]
	switch command {
	case "/get_location_update":
		return getFreshLocation(tripCode, chatId)
	}
	return nil, nil
}

func getFreshLocation(tripCode string, chatId int64) ([]*telego.SendMessageParams, error) {
	return nil, nil
}
