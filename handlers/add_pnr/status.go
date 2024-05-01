package addpnr

type Init struct{}

type SaveTripCode struct {
	TripCode string
}

type SaveTripCodePnr struct {
	Pnr      string
	TripCode string
}
type SaveTripCodePnrProvider struct {
	Pnr             string
	TripCode        string
	ServiceProvider string
}

type SubmittedForVerification struct {
	Pnr             string
	TripCode        string
	ServiceProvider string
}
