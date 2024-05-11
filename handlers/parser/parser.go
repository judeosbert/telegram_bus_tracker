package parser

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"
)

type ticketInfo struct {
	ServiceProvider string
	BusNumber       string
	Doj             time.Time
}

type ticketParser struct {
}

// ParseTicket implements TicketParser.
func (t *ticketParser) ParseTicket(ticket string) (*ticketInfo, error) {
	res, err := parseKeralaTicket(ticket)
	if err == nil {
		return res, nil
	}
	res, err = parseKarnatakaTicket(ticket)
	if err == nil {
		return res, nil
	}
	return nil, errors.New("invalid ticket")

}

type TicketParser interface {
	ParseTicket(ticket string) (*ticketInfo, error)
}

func NewTicketParser() TicketParser {
	return &ticketParser{}
}

//Sample ticket
/**
KSRTC Bus PNR: J101083761, DOJ: 06-May-2024 23:58, 2201KZKBNG, Bus No : KA-40 F-1196, Crew Mobile No: 9164455494. Happy Journey. From MKSRTC
**/
func parseKarnatakaTicket(ticket string) (*ticketInfo, error) {
	if !strings.Contains(ticket, "From MKSRTC") {
		return nil, errors.ErrUnsupported
	}
	t := &ticketInfo{}
	r, err := regexp.Compile("(DOJ: )([A-Z0-9a-z- :]+)")
	if err != nil {
		log.Println("Error in regex", err)
		return nil, err
	}
	dojString := strings.Replace(string(r.Find([]byte(ticket))), "DOJ: ", "", 1)
	doj, err := time.Parse("02-Jan-2006 15:04", dojString)
	if err != nil {
		log.Println("Error in parsing date", err)
		return nil, err
	}
	t.Doj = doj

	r, err = regexp.Compile("(Bus No : )([A-Z0-9a-z- ]+)")
	if err != nil {
		log.Println("Error in parsing bus no regex ", err)
		return nil, err
	}
	busNoString := strings.Split(string(r.Find([]byte(ticket))), ":")[1]
	busNo := strings.TrimSpace(busNoString)
	if len(busNo) == 0 {
		log.Println("Error in parsing bus no")
		return nil, errors.New("Bus no not found")
	}
	t.BusNumber = busNo
	t.ServiceProvider = "Karnataka RTC"

	return t, nil
}

//Sample ticket
/**
Dear Jude, Your Bus NO:KS113, FOR Kozhikode-Mananthavady, Operator NAME KOZHIKODE DEPOT PNR NO:4792692,Pickup:Kozhikode, Crew NAME :LIJO, Crew NO:9947076697. Plz reach 15 mins BEFORE your journey TIME:09:29 PM. - Happy Journey - Kerala RTC
**/
func parseKeralaTicket(ticket string) (*ticketInfo, error) {
	if !strings.Contains(ticket, "Kerala RTC") {
		return nil, errors.New("Invalid ticket for Kerala RTC")
	}
	t := &ticketInfo{}
	r, err := regexp.Compile("(Your Bus NO:)([A-Z0-9a-z- ]+)")
	if err != nil {
		log.Println("Error in parsing bus no regex ", err)
		return nil, err
	}
	busNoString := string(r.Find([]byte(ticket)))
	busNoString = strings.Split(busNoString, ":")[1]
	busNo := strings.TrimSpace(busNoString)
	if len(busNo) == 0 {
		log.Println("Error in parsing bus no")
		return nil, errors.New("Bus no not found")
	}
	t.BusNumber = busNo
	t.ServiceProvider = "Kerala SRTC"
	return t, nil
}
