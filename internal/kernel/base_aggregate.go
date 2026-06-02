package kernel

import "time"

type BaseAggregate struct {
	id                string
	version           int
	createdAt         time.Time
	uncommittedEvents []DomainEvent
}

type HasDomainEvents interface {
	Id() string
	Version() int
	IncrementVersion()
	UncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
}

func (b *BaseAggregate) Id() string {
	return b.id
}

func (b *BaseAggregate) Version() int {
	return b.version
}

func (b *BaseAggregate) IncrementVersion() {
	b.version++
}

func (b *BaseAggregate) SetId(id string) {
	b.id = id
}

func (b *BaseAggregate) CreatedAt() time.Time {
	return b.createdAt
}

func (b *BaseAggregate) SetCreatedAt(value time.Time) {
	b.createdAt = value
}

func (b *BaseAggregate) ApplyEvent(event DomainEvent, applyFunc func(DomainEvent)) {
	applyFunc(event)
	b.uncommittedEvents = append(b.uncommittedEvents, event)
}

func (b *BaseAggregate) UncommittedEvents() []DomainEvent {
	return b.uncommittedEvents
}

func (b *BaseAggregate) ClearUncommittedEvents() {
	b.uncommittedEvents = nil
}

func (b *BaseAggregate) LoadFromHistory(events []DomainEvent, applyFunc func(DomainEvent)) {
	for _, e := range events {
		applyFunc(e)
		b.IncrementVersion()
	}
}
