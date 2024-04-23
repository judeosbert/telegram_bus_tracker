package addpnr
type Init struct{}

type SavePnr struct {
	Pnr string
}

type SavePnrTripCode struct {
	Pnr string
	TripCode string
}
type SavePnrTripCodeProvider struct {
	Pnr string
	TripCode string
	ServiceProvider string
}

type SubmittedForVerification struct {
	Pnr string
	TripCode string
	ServiceProvider string
}