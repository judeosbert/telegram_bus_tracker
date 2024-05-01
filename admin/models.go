package admin

type NewTripInfo struct {
	ServiceProvider string
	TripCode        string
	Pnr             string
}

type TripStateVerifiedCallback struct {
	Type     string `json:"type"`
	Pnr      string `json:"pnr"`
	TripCode string `json:"trip_code"`
	Status   string `json:"status"`
}

type TripStateValidation struct {
	Pnr      string `json:"pnr"`
	TripCode string `json:"trip_code"`
	Status   string `json:"status"`
}

func NewStateTripVerified(pnr string, tripCode string) TripStateValidation {
	return TripStateValidation{
		Pnr:      pnr,
		TripCode: tripCode,
		Status:   STATUS_VERIFIED_TRIP_VERIFICATION,
	}
}
func NewStateTripRejected(pnr string, tripCode string) TripStateValidation {
	return TripStateValidation{
		Pnr:      pnr,
		TripCode: tripCode,
		Status:   STATUS_REJECTED_TRIP_VERIFICATION,
	}
}

type TripStateVerifiedWithLink struct {
	Pnr      string `json:"pnr"`
	TripCode string `json:"trip_code"`
	InviteLink string `json:"link"`
}

var TYPE_TRIP_VERIFICATION = "trip-verification"
var STATUS_VERIFIED_TRIP_VERIFICATION = "verified"
var STATUS_REJECTED_TRIP_VERIFICATION = "rejected"

func NewTripRejectedStateCallback(pnr string, tripcode string) TripStateVerifiedCallback {
	return TripStateVerifiedCallback{
		Type:     TYPE_TRIP_VERIFICATION,
		Pnr:      pnr,
		TripCode: tripcode,
		Status:   STATUS_REJECTED_TRIP_VERIFICATION,
	}
}
func NewTripVerifiedStateCallback(pnr string, tripcode string) TripStateVerifiedCallback {
	return TripStateVerifiedCallback{
		Type:     TYPE_TRIP_VERIFICATION,
		Pnr:      pnr,
		TripCode: tripcode,
		Status:   STATUS_VERIFIED_TRIP_VERIFICATION,
	}
}
