package journey

type CreateJourneyLogCommand struct {
	CardNo            string
	PredecessorCardNo string
}

type StartJourneyCommand struct {
	CardNo       string
	StartStation string
}

type EndJourneyCommand struct {
	CardNo     string
	EndStation string
}
