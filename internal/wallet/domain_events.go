package wallet

import "time"

const (
	EventBalanceAdded = "BalanceAddedDomainEvent"
	EventFareDeducted = "FareDeductedDomainEvent"
	EventCardCreated  = "CardCreatedDomainEvent"
)

type CardCreatedDomainEvent struct {
	Id        string
	CardNo    string
	CreatedAt time.Time
	ValidTo   time.Time
}

type BalanceAddedDomainEvent struct {
	BalanceAdded int
}

type FareDeductedDomainEvent struct {
	DeductedFare int
}

func (e CardCreatedDomainEvent) EventName() string {
	return EventCardCreated
}

func (e BalanceAddedDomainEvent) EventName() string {
	return EventBalanceAdded
}

func (e FareDeductedDomainEvent) EventName() string {
	return EventFareDeducted
}
