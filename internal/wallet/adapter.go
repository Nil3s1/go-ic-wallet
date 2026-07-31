package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type PaymentAdapter struct {
	repo  CardProjectionRepository
	store kernel.EventStore[*Card]
}

func NewPaymentAdapter(repo CardProjectionRepository, store kernel.EventStore[*Card]) *PaymentAdapter {
	return &PaymentAdapter{
		repo:  repo,
		store: store,
	}
}

func (p *PaymentAdapter) HasSufficientBalance(cardNo string, amount int) (bool, error) {
	query := HasSufficientBalanceQuery{
		CardNo: cardNo,
		Amount: amount,
	}
	handler := NewHasSufficientBalanceQueryHandler(p.repo)
	return handler.Handle(query)
}

func (p *PaymentAdapter) AuthorizePayment(ctx context.Context, cardNo string, referenceId string, amount int) error {
	command := ApplyPaymentCommand{
		CardNo:      cardNo,
		Amount:      amount,
		ReferenceId: referenceId,
	}
	handler := NewApplyPaymentCommandHandler(p.store)
	return handler.Handle(ctx, command)
}
