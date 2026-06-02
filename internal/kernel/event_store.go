package kernel

import "context"

type EventStore[T HasDomainEvents] interface {
	Load(ctx context.Context, id string) (T, error)
	Save(ctx context.Context, aggregate T) error
}
