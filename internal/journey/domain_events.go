package journey

import "time"

const (
	EventJourneyLogCreated = "JourneyLogCreatedEvent"
	EventJourneyStarted    = "JourneyStartedDomainEvent"
	EventJourneyEnded      = "JourneyEndedDomainEvent"
)

type JourneyLogCreatedDomainEvent struct {
	PredecessorCardNo string
	CardNo            string
	CreatedAt         time.Time
}

type JourneyStartedDomainEvent struct {
	StartStation string
	StartTime    time.Time
}

type JourneyEndedDomainEvent struct {
	EndStation        string
	EndTime           time.Time
	DistanceTravelled int
	Fare              int
}

func (e JourneyLogCreatedDomainEvent) EventName() string {
	return EventJourneyLogCreated
}

func (e JourneyStartedDomainEvent) EventName() string {
	return EventJourneyStarted
}

func (e JourneyEndedDomainEvent) EventName() string {
	return EventJourneyEnded
}
