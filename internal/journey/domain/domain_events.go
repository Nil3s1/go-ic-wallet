package domain

import "time"

const (
	EventJourneyLogCreatedV1 = "JourneyLogCreatedEvent.V1"
	EventJourneyStartedV1    = "JourneyStartedDomainEvent.V1"
	EventJourneyEndedV1      = "JourneyEndedDomainEvent.V1"
)

type JourneyLogCreatedDomainEventV1 struct {
	MediaId   string
	CreatedAt time.Time
}

type JourneyStartedDomainEventV1 struct {
	StartStation       string
	StartTime          time.Time
	JourneyReferenceId string
}

type JourneyEndedDomainEventV1 struct {
	EndStation        string
	EndTime           time.Time
	DistanceTravelled uint
	Fare              uint
}

func (e JourneyLogCreatedDomainEventV1) EventName() string {
	return EventJourneyLogCreatedV1
}

func (e JourneyStartedDomainEventV1) EventName() string {
	return EventJourneyStartedV1
}

func (e JourneyEndedDomainEventV1) EventName() string {
	return EventJourneyEndedV1
}
