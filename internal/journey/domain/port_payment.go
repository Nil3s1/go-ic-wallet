package domain

import "context"

type PaymentPort interface {
	HasSufficientBalance(ctx context.Context, mediaId string, amount int) (bool, error)
	AuthorizePayment(ctx context.Context, mediaId string, referenceId string, amount int) error
}
