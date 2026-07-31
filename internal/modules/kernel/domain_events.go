package kernel

type DomainEvent interface {
	EventName() string
}
