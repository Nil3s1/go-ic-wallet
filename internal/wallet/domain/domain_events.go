package domain

import "time"

const (
	EventBalanceAdded = "BalanceAddedDomainEvent"
	EventApplyPayment = "ApplyPaymentDomainEvent"
	EventCardCreated  = "CardCreatedDomainEvent"
)

type CardCreatedDomainEvent struct {
	CardNo         string
	InitialBalance uint
	CreatedAt      time.Time
	ValidTo        time.Time
}

type BalanceAddedDomainEvent struct {
	BalanceAdded uint
	ReferenceId  string
}

type ApplyPaymentDomainEvent struct {
	Amount      uint
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
