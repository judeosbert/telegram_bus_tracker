package addpnr
type Init struct{}

type SavePnr struct {
	Pnr string
}

type SavePnrProvider struct {
	Pnr string
	ServiceProvider string
}

type SubmittedForVerification struct {
	Pnr string
	ServiceProvider string
}