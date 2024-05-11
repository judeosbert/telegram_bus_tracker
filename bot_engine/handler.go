package botengine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/judeosbert/bus_tracker_bot/admin"
	addpnr "github.com/judeosbert/bus_tracker_bot/handlers/add_pnr"
	"github.com/judeosbert/bus_tracker_bot/handlers/parser"
	"github.com/judeosbert/bus_tracker_bot/handlers/start"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/judeosbert/bus_tracker_bot/utils"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleAfterServiceProvideSelected(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("No previous state")
	}
	switch prevState.(type) {
	case addpnr.Init:
		break
	default:
		return nil, errors.New("Invalid previous state")
	}

	msg := message.Text
	switch msg {
	case "Kerala SRTC":
		saver.SetUserState(chatId, addpnr.ServiceProviderSet{Provider: "Kerala SRTC"})
	case "Karnataka SRTC":
		saver.SetUserState(chatId, addpnr.ServiceProviderSet{Provider: "Karnataka SRTC"})
	default:

		return []*telego.SendMessageParams{
			{
				ChatID: tu.ID(chatId),
				Text:   "Sorry invalid service provider. Try again. Should be Kerala SRTC or Karnataka SRTC",
			},
		}, nil
	}

	return []*telego.SendMessageParams{
		{
			ChatID: tu.ID(chatId),
			Text:   "Okay, send the bus number.",
		},
	}, nil
}

func HandleAfterBusNoSent(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("No previous state")
	}
	switch prevState := prevState.(type) {
	case addpnr.ServiceProviderSet:
		msg := message.Text
		if len(msg) == 0 {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Invalid bus number. Try again",
				},
			}, nil
		}
		msg = strings.ToUpper(msg)
		saver.SetUserState(chatId, addpnr.ServiceProviderBusNoSet{BusNo: msg, Provider: prevState.Provider})
		return []*telego.SendMessageParams{
			{
				ChatID: tu.ID(chatId),
				Text:   "Okay, when is the bus starting? Send the date in the format dd/mm/yyyy",
			},
		}, nil
	default:
		return nil, errors.New("Invalid previous state")
	}
}

func HandleAfterDojSent(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("No previous state")
	}
	switch prevState := prevState.(type) {
	case addpnr.ServiceProviderBusNoSet:
		msg := message.Text
		if len(msg) == 0 {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Invalid date. Try again",
				},
			}, nil
		}
		doj, err := time.Parse("02/01/2006", msg)
		today, _ := time.Parse("02/01/2006", time.Now().AddDate(0, 0, -1).Format("02/01/2006"))
		if doj.Before(today) {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Invalid date. Date should be in the future.",
				},
			}, nil
		}
		if err != nil {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   err.Error(),
				},
			}, nil
		}
		return addOrCreateGroup(adminUtils, saver, prevState.BusNo, doj, chatId)

	default:
		return nil, errors.New("invalid previous state")
	}
}

func addOrCreateGroup(adminUtils admin.AdminUtils, saver state.Saver, busNo string, doj time.Time, chatId int64) ([]*telego.SendMessageParams, error) {
	//Check if a group exists for this trip.
	_, err := saver.GetTripGroup(utils.TripHash(busNo, doj))
	if err == nil {
		saver.AddTripObserver(utils.TripHash(busNo, doj), chatId)
		adminUtils.AddToTripGroup(admin.NewTripInfo{
			Doj:   doj,
			BusNo: busNo,
		}, tu.ID(chatId))
		return []*telego.SendMessageParams{}, nil
	}

	saver.AddTripObserver(utils.TripHash(busNo, doj), chatId)
	adminUtils.SubmitForNewGroup(admin.NewTripInfo{
		Doj:   doj,
		BusNo: busNo,
	})
	saver.RemoveUserState(chatId)
	return []*telego.SendMessageParams{
		{
			ChatID: tu.ID(chatId),
			Text:   "Okay, you will be notified when the group is created.Then join the group for updates.",
		},
	}, nil
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
	// saver.SetTripGroup(tripCode, chatId)

	//supported commands get update.
	return nil, nil
}

func AdminHandleNewGroupHashtag(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	if message.From.ID != 885727411 {
		return nil, errors.New("not an admin message")
	}
	if message.Entities == nil {
		return nil, errors.New("no entities")
	}
	var hasHashtag = false
	for i := 0; i < len(message.Entities); i++ {
		e := message.Entities[i]
		if e.Type == "hashtag" {
			hasHashtag = true
			break
		}
	}
	if !hasHashtag {
		return nil, errors.New("not a hashtag message")
	}
	text := message.Text
	if !strings.Contains(text, "#newgroup") {
		return nil, errors.New("invalid hashtag")
	}
	chatId := message.Chat.ID
	adminUtils.OnNewGroup(chatId)
	id := message.MessageID
	adminUtils.DeleteMessages(chatId, []int{id})

	return []*telego.SendMessageParams{
		{
			ChatID: tu.ID(885727411),
			Text:   "Thanks! New group used.",
		},
	}, nil
}

func HandleAfterTicketSent(message telego.Message, saver state.Saver, adminUtils admin.AdminUtils) ([]*telego.SendMessageParams, error) {
	chatId := message.Chat.ID
	prevState, err := saver.GetUserState(chatId)
	if err != nil {
		return nil, errors.New("No previous state")
	}
	switch prevState.(type) {
	case start.RequestTicket:
		msg := message.Text
		if len(msg) == 0 {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Invalid ticket. Try again",
				},
			}, nil
		}

		t, err := parser.NewTicketParser().ParseTicket(msg)
		if err != nil {
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Invalid ticket. Try again",
				},
			}, nil
		}
		if t.ServiceProvider == "Kerala SRTC" {
			saver.SetUserState(chatId, addpnr.ServiceProviderBusNoSet{
				BusNo:    t.BusNumber,
				Provider: t.ServiceProvider,
			})
			return []*telego.SendMessageParams{
				{
					ChatID: tu.ID(chatId),
					Text:   "Okay, which date is the bus starting? Send dd/mm/yyyy format",
				},
			}, nil
		} else {
			saver.SetUserState(chatId, addpnr.ServiceProviderBusNoDojSet{
				BusNo:    t.BusNumber,
				Provider: t.ServiceProvider,
				Doj:      t.Doj,
			})
			return addOrCreateGroup(adminUtils, saver, t.BusNumber, t.Doj, chatId)
		}

	default:
		return nil, errors.New("Invalid previous state")
	}
}
