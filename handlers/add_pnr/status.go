package addpnr

import "time"

type Init struct{}

type ServiceProviderSet struct {
	Provider string
}

type ServiceProviderBusNoSet struct {
	BusNo string
	Provider string
}

type ServiceProviderBusNoDojSet struct {
	BusNo string
	Provider string
	Doj time.Time
}

