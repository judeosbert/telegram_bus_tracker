package admin

import (
	"fmt"
	"log"
	"time"

	"github.com/judeosbert/bus_tracker_bot/handlers"
	"github.com/judeosbert/bus_tracker_bot/state"
	"github.com/judeosbert/bus_tracker_bot/utils"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type adminUtils struct {
	groupAssigner GroupAssigner
	admin         telego.ChatID
	bot           *telego.Bot
	saver         state.Saver
}

// DeleteMessages implements AdminUtils.
func (a *adminUtils) DeleteMessages(chatId int64, messageId []int) {
	a.bot.DeleteMessages(&telego.DeleteMessagesParams{
		ChatID:     tu.ID(chatId),
		MessageIDs: messageId,
	})
}

// OnNewGroup implements AdminUtils.
func (a *adminUtils) OnNewGroup(link int64) {
	go func() {
		a.groupAssigner.OnNewGroup(link)
	}()
}

// AddToTripGroup implements AdminUtils.
func (a *adminUtils) AddToTripGroup(trip NewTripInfo, chatID telego.ChatID) {
	grp, err := a.saver.GetTripGroup(utils.TripHash(trip.BusNo, trip.Doj))
	if err != nil {
		handlers.SendMessage(a.bot, chatID.ID, "Error getting group")
		return
	}
	link, err := a.bot.CreateChatInviteLink(&telego.CreateChatInviteLinkParams{
		ChatID:             tu.ID(grp),
		Name:               "Welcome to the group",
		ExpireDate:         time.Now().AddDate(0, 0, 1).Unix(),
		MemberLimit:        1,
		CreatesJoinRequest: false,
	})
	if err != nil {
		handlers.SendMessage(a.bot, chatID.ID, "Error creating group link. Try again later.")
		return
	}
	handlers.SendMessage(a.bot, chatID.ID, fmt.Sprintf("Join this group for updates on your on %s bus %s\n%s", trip.Doj.Format("02/06/2024"), trip.BusNo, link.InviteLink))
}

// SubmitForNewGroup implements AdminUtils.
func (a *adminUtils) SubmitForNewGroup(trip NewTripInfo) {
	// Send message to admin
	a.groupAssigner.AssignGroup(trip)
	params := telego.SendMessageParams{
		ChatID: a.admin,
		Text:   "Need new group!",
	}
	a.sendMessage(params)
}

func (a *adminUtils) sendMessage(params telego.SendMessageParams) error {
	sentMessage, err := a.bot.SendMessage(
		&params,
	)
	if err != nil {
		return err
	}

	log.Printf("Admin:Sent message %+v\n", sentMessage)
	return nil
}

type AdminUtils interface {
	AddToTripGroup(trip NewTripInfo, chatID telego.ChatID)
	SubmitForNewGroup(trip NewTripInfo)
	OnNewGroup(link int64)
	DeleteMessages(chatId int64, messageId []int)
}

func NewAdminUtils(bot *telego.Bot, saver state.Saver) AdminUtils {
	groupAssigner := NewGroupAssigner()
	adminUtils := &adminUtils{
		saver:         saver,
		groupAssigner: groupAssigner,
		bot:           bot,
		admin:         tu.ID(885727411),
	}
	go func() {
		for group := range groupAssigner.ResultChan() {
			saver.SetTripGroup(utils.TripHash(group.NewTripInfo.BusNo, group.NewTripInfo.Doj), group.GroupId)
			bot.SetChatTitle(&telego.SetChatTitleParams{
				ChatID: telego.ChatID{
					ID: group.GroupId,
				},
				Title: fmt.Sprintf("%s on %s", group.NewTripInfo.BusNo, group.NewTripInfo.Doj.Format("02/01/2006")),
			})
			msg, _ := bot.SendMessage(&telego.SendMessageParams{
				ChatID: telego.ChatID{
					ID: group.GroupId,
				},
				Text: "This is a public group. Be friendly in the chats. No spamming or sharing of inappropriate content. We are not responsible for any damages caused by the group members.",
			})
			bot.PinChatMessage(&telego.PinChatMessageParams{
				ChatID: telego.ChatID{
					ID: group.GroupId,
				},
				MessageID:           msg.GetMessageID(),
				DisableNotification: false,
			})
			obs := adminUtils.saver.GetTripObservers(utils.TripHash(group.NewTripInfo.BusNo, group.NewTripInfo.Doj))
			for _, ob := range obs {
				adminUtils.AddToTripGroup(NewTripInfo{
					Doj:   group.Doj,
					BusNo: group.BusNo,
				}, tu.ID(ob))
			}
		}
	}()
	return adminUtils
}
