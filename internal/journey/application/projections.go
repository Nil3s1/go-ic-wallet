package application

import (
	"context"
	"time"
)

type JourneyLogProjection struct {
	MediaId string
	// LastJourneyReferenceId is kept after the journey ended so replays can still resolve it.
	LastJourneyReferenceId string
	IsOnJourney            bool
	LastStation            string
	UpdatedAt              time.Time
}

type JourneyLogProjectionRepository interface {
	GetByKey(ctx context.Context, mediaId string) (JourneyLogProjection, error)
}

const (
	JourneyStatusStarted = "Started"
	JourneyStatusEnded   = "Ended"
)

type JourneyProjection struct {
	JourneyReferenceId string
	MediaId            string
	StartStation       string
	StartTime          time.Time
	EndStation         string
	EndTime            *time.Time // nil while the journey is still running
	DistanceTravelled  uint
	Fare               uint //Currency in cents
	Status             string
}

type JourneyProjectionRepository interface {
	GetByKey(ctx context.Context, journeyReferenceId string) (JourneyProjection, error)
	ListByMedia(ctx context.Context, mediaId string, limit int, offset int) ([]JourneyProjection, error)
}
