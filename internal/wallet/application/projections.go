package application

import (
	"context"
	"time"
)

type CardProjection struct {
	CardNo         string
	ValidTo        time.Time
	CurrentBalance uint //Currency in cents
}

type CardProjectionRepository interface {
	GetCard(ctx context.Context, cardNo string) (CardProjection, error)
	HasSufficientBalance(ctx context.Context, cardNo string, amount uint) (bool, error)

	Update(ctx context.Context, card CardProjection) error
	Add(ctx context.Context, card CardProjection) error
}
