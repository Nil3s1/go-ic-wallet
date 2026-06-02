package kernel

import "time"

type BaseAggregate struct {
	id                string
	createdAt         time.Time
	uncommittedEvents []DomainEvent
}

type Persistable interface {
	Id() string
	UncommittedEvents() []DomainEvent
	ClearUncommittedEvents()
}

func (b *BaseAggregate) Id() string {
	return b.id
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

func (b *BaseAggregate) AddEvent(event DomainEvent) {
	b.uncommittedEvents = append(b.uncommittedEvents, event)
}

func (b *BaseAggregate) UncommittedEvents() []DomainEvent {
	return b.uncommittedEvents
}

func (b *BaseAggregate) ClearUncommittedEvents() {
	b.uncommittedEvents = nil
}
