package journey

import (
	"errors"
	"time"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type JourneyLog struct {
	kernel.BaseAggregate

	predecessorCardNo string
	isOnJourney       bool
	startStation      string
	startTime         time.Time
	endStation        string
	endTime           time.Time
	distanceTravelled int
	fare              int
}

func NewJourneyLog(cardNo string) (*JourneyLog, error) {
	return NewJourneyLogWithPredecessor(cardNo, "")
}

func NewJourneyLogWithPredecessor(cardNo string, predecessorCardNo string) (*JourneyLog, error) {
	if cardNo == "" {
		return nil, errors.New("cannot create journey log: cardNo is empty")
	}

	createdAt := time.Now().UTC()

	event := JourneyLogCreatedDomainEvent{
		CardNo:            cardNo,
		CreatedAt:         createdAt,
		PredecessorCardNo: predecessorCardNo,
	}

	jl := &JourneyLog{}
	jl.ApplyEvent(event, jl.applyEventFunction)

	return jl, nil
}

func Rehydrate(events []kernel.DomainEvent) *JourneyLog {
	jl := &JourneyLog{}
	jl.LoadFromHistory(events, jl.applyEventFunction)

	return jl
}

func (jl *JourneyLog) PredecessorCardNo() string {
	return jl.predecessorCardNo
}

func (jl *JourneyLog) CardNo() string {
	return jl.Id()
}

func (jl *JourneyLog) IsOnJourney() bool {
	return jl.isOnJourney
}

func (jl *JourneyLog) StartStation() string {
	return jl.startStation
}

func (jl *JourneyLog) StartTime() time.Time {
	return jl.startTime
}

func (jl *JourneyLog) EndStation() string {
	return jl.endStation
}

func (jl *JourneyLog) EndTime() time.Time {
	return jl.endTime
}

func (jl *JourneyLog) DistanceTravelled() int {
	return jl.distanceTravelled
}

func (jl *JourneyLog) Fare() int {
	return jl.fare
}

func (jl *JourneyLog) StartJourney(startStation string) error {
	if startStation == "" {
		return errors.New("Please Provide a StartStation!")
	}
	if jl.isOnJourney {
		return errors.New("There is already a Started Journey for this card. Please TapOut First!")
	}

	startTime := time.Now().UTC()

	event := JourneyStartedDomainEvent{
		StartStation: startStation,
		StartTime:    startTime,
	}

	jl.ApplyEvent(event, jl.applyEventFunction)

	return nil
}

func (jl *JourneyLog) EndJourney(endStation string, cf CalculatedFare) error {
	if endStation == "" {
		return errors.New("please provide a valid endStation")
	}
	if !jl.isOnJourney {
		return errors.New("cannot end journey: no journey started for this card")
	}
	if cf.Distance() < 0 || cf.Fare() < 0 {
		return errors.New("distance and fare must be positive")
	}

	event := JourneyEndedDomainEvent{
		EndStation:        endStation,
		EndTime:           time.Now().UTC(),
		DistanceTravelled: cf.Distance(),
		Fare:              cf.Fare(),
	}

	jl.ApplyEvent(event, jl.applyEventFunction)
	return nil
}

func (jl *JourneyLog) applyEventFunction(event kernel.DomainEvent) {
	switch e := event.(type) {
	case JourneyLogCreatedDomainEvent:
		jl.SetId(e.CardNo)
		jl.SetCreatedAt(e.CreatedAt)
		jl.predecessorCardNo = e.PredecessorCardNo
		jl.isOnJourney = false
	case JourneyStartedDomainEvent:
		jl.isOnJourney = true
		jl.startStation = e.StartStation
		jl.startTime = e.StartTime
		jl.endStation = ""
		jl.endTime = time.Time{}
		jl.distanceTravelled = 0
		jl.fare = 0
	case JourneyEndedDomainEvent:
		jl.isOnJourney = false
		jl.endStation = e.EndStation
		jl.endTime = e.EndTime
		jl.distanceTravelled = e.DistanceTravelled
		jl.fare = e.Fare
	}
}
