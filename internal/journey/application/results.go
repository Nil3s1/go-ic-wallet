package application

type StartJourneyResult struct {
	MediaId     string
	IsOnJourney bool
}

type EndJourneyResult struct {
	MediaId           string
	IsOnJourney       bool
	DistanceTravelled uint
	Fare              uint
}
