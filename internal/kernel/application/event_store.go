package application

import (
	"context"

	kernelDomain "github.com/Nil3s1/go-ic-wallet/internal/kernel/domain"
)

type EventStore[T kernelDomain.HasDomainEvents] interface {
	Exists(ctx context.Context, id string) (bool, error)
	Load(ctx context.Context, id string) (T, error)
	Save(ctx context.Context, aggregate T) error
}
