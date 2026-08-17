package infrastructure

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type PaymentAdapter struct {
	repo  domain.CardProjectionRepository
	store kernel.EventStore[*domain.Card]
}

func NewPaymentAdapter(repo domain.CardProjectionRepository, store kernel.EventStore[*domain.Card]) *PaymentAdapter {
	return &PaymentAdapter{
		repo:  repo,
		store: store,
	}
}

func (p *PaymentAdapter) HasSufficientBalance(cardNo string, amount int) (bool, error) {
	query := application.HasSufficientBalanceQuery{
		CardNo: cardNo,
		Amount: amount,
	}
	handler := application.NewHasSufficientBalanceQueryHandler(p.repo)
	return handler.Handle(query)
}

func (p *PaymentAdapter) AuthorizePayment(ctx context.Context, cardNo string, referenceId string, amount int) error {
	command := application.ApplyPaymentCommand{
		CardNo:      cardNo,
		Amount:      amount,
		ReferenceId: referenceId,
	}
	handler := application.NewApplyPaymentCommandHandler(p.store)
	return handler.Handle(ctx, command)
}
