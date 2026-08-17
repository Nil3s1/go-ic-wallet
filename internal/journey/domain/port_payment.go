package domain

import "context"

type PaymentPort interface {
	HasSufficientBalance(ctx context.Context, mediaId string, amount uint) (bool, error)
	AuthorizePayment(ctx context.Context, mediaId string, referenceId string, amount uint) error
}
