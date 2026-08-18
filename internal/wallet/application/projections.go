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
	GetByKey(ctx context.Context, cardNo string) (CardProjection, error)
}

const (
	BookingDirectionCredit = "Credit"
	BookingDirectionDebit  = "Debit"
)

type BookingProjection struct {
	CardNo       string
	ReferenceId  string // opaque idempotency token, the wallet does not interpret it
	Amount       uint   //Currency in cents
	Direction    string
	BookedAt     time.Time
	BalanceAfter uint //Currency in cents
}

type BookingProjectionRepository interface {
	GetByReference(ctx context.Context, cardNo string, referenceId string) (BookingProjection, error)
	ListByCard(ctx context.Context, cardNo string, limit int, offset int) ([]BookingProjection, error)
}
