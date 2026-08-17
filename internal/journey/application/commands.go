package application

type StartJourneyCommand struct {
	MediaId      string
	StartStation string
}

type EndJourneyCommand struct {
	MediaId    string
	EndStation string
}
