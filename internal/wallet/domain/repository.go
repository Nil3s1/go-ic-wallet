package domain

import "context"

type CardProjectionRepository interface {
	GetCard(ctx context.Context, cardNo string) (CardProjection, error)
	HasSufficientBalance(ctx context.Context, cardNo string, amount uint) (bool, error)

	Update(ctx context.Context, card CardProjection) error
	Add(ctx context.Context, card CardProjection) error
}
