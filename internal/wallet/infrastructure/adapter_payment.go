package infrastructure

import (
	"context"

	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type PaymentAdapter struct {
	repo                        domain.CardProjectionRepository
	store                       kernelApplication.EventStore[*domain.Card]
	hasSufficientBalanceHandler *application.HasSufficientBalanceQueryHandler
	applyPaymentHandler         *application.ApplyPaymentCommandHandler
}

func NewPaymentAdapter(
	repo domain.CardProjectionRepository,
	store kernelApplication.EventStore[*domain.Card],
	hasSufficientBalanceHandler *application.HasSufficientBalanceQueryHandler,
	applyPaymentHandler *application.ApplyPaymentCommandHandler) *PaymentAdapter {
	return &PaymentAdapter{
		repo:                        repo,
		store:                       store,
		hasSufficientBalanceHandler: hasSufficientBalanceHandler,
		applyPaymentHandler:         applyPaymentHandler,
	}
}

func (p *PaymentAdapter) HasSufficientBalance(ctx context.Context, cardNo string, amount uint) (bool, error) {
	query := application.HasSufficientBalanceQuery{
		CardNo: cardNo,
		Amount: amount,
	}
	return p.hasSufficientBalanceHandler.Handle(ctx, query)
}

func (p *PaymentAdapter) AuthorizePayment(ctx context.Context, cardNo string, referenceId string, amount uint) error {
	command := application.ApplyPaymentCommand{
		CardNo:      cardNo,
		Amount:      uint(amount),
		ReferenceId: referenceId,
	}
	return p.applyPaymentHandler.Handle(ctx, command)
}
