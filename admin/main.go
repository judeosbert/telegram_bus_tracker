package admin

import (
	"fmt"
	"log"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type adminUtils struct {
	admin telego.ChatID
	bot   *telego.Bot
}

// SendForVerification implements AdminUtils.
func (a *adminUtils) SendForVerification(trip NewTripInfo) {
	var webUrl string
	if trip.ServiceProvider == "Kerala SRTC" {
		webUrl = fmt.Sprintf("https://onlineksrtcswift.com/status/booking/%s?u=true", trip.Pnr)
	} else {
		webUrl = fmt.Sprintf("https://www.ksrtc.in/oprs-web/tickets/status/serviceDetails.do?pnrPrefixWithTktNo=%s&serviceCode=0&txtDepartureDate=undefined", trip.Pnr)
	}
	msg := telego.SendMessageParams{
		ChatID: a.admin,
		Text:   fmt.Sprintf("New Trip For Verification:\n PNR:%s\n TripCode:%s\n Service Provider:%s ", trip.Pnr, trip.TripCode,trip.ServiceProvider),
		ReplyMarkup: &telego.InlineKeyboardMarkup{
			InlineKeyboard: [][]telego.InlineKeyboardButton{
				{
					telego.InlineKeyboardButton{
						Text:         "Open Service Provider Page",
						URL:          webUrl,
					},
				},
				{
					telego.InlineKeyboardButton{
						Text:         "Mark as Verified",
						CallbackData: fmt.Sprintf("verified-%s", trip.Pnr),
					},
					telego.InlineKeyboardButton{
						Text:         "Mark as Rejected",
						CallbackData: fmt.Sprintf("rejected-%s", trip.Pnr),
					},
				},
			},
		},
	}
	a.sendMessage(msg)
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
	SendForVerification(trip NewTripInfo)
}

func NewAdminUtils(bot *telego.Bot) AdminUtils {
	return &adminUtils{
		bot:   bot,
		admin: tu.ID(885727411),
	}
}
