package domain

import "time"

const (
	EventBalanceAddedV1 = "BalanceAddedDomainEvent.V1"
	EventApplyPaymentV1 = "ApplyPaymentDomainEvent.V1"
	EventCardCreatedV1  = "CardCreatedDomainEvent.V1"
)

type CardCreatedDomainEventV1 struct {
	CardNo         string
	InitialBalance uint
	CreatedAt      time.Time
	ValidTo        time.Time
}

type BalanceAddedDomainEventV1 struct {
	BalanceAdded uint
	ReferenceId  string
	OccurredAt   time.Time
}

type ApplyPaymentDomainEventV1 struct {
	Amount      uint
	ReferenceId string
	OccurredAt  time.Time
}

func (e CardCreatedDomainEventV1) EventName() string {
	return EventCardCreatedV1
}

func (e BalanceAddedDomainEventV1) EventName() string {
	return EventBalanceAddedV1
}

func (e ApplyPaymentDomainEventV1) EventName() string {
	return EventApplyPaymentV1
}
