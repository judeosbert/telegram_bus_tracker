package admin

type NewTripInfo struct{
	ServiceProvider string
	TripCode string
	Pnr string
}

type TripStateVerified struct {
	Pnr string
	TripCode string
}

type TripStateRejected struct {
	Pnr string
	Reason string
}