package domain

import "time"

const (
	EventBalanceAdded = "BalanceAddedDomainEvent"
	EventApplyPayment = "ApplyPaymentDomainEvent"
	EventCardCreated  = "CardCreatedDomainEvent"
)

type CardCreatedDomainEvent struct {
	CardNo         string
	InitialBalance int
	CreatedAt      time.Time
	ValidTo        time.Time
}

type BalanceAddedDomainEvent struct {
	BalanceAdded int
	ReferenceId  string
}

type ApplyPaymentDomainEvent struct {
	Amount      int
	ReferenceId string
}

func (e CardCreatedDomainEvent) EventName() string {
	return EventCardCreated
}

func (e BalanceAddedDomainEvent) EventName() string {
	return EventBalanceAdded
}

func (e ApplyPaymentDomainEvent) EventName() string {
	return EventApplyPayment
}
