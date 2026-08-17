package application

import (
	"context"

	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type ApplyPaymentCommandHandler struct {
	store kernelApplication.EventStore[*domain.Card]
}

func NewApplyPaymentCommandHandler(store kernelApplication.EventStore[*domain.Card]) *ApplyPaymentCommandHandler {
	return &ApplyPaymentCommandHandler{
		store: store,
	}
}

func (h *ApplyPaymentCommandHandler) Handle(ctx context.Context, cmd ApplyPaymentCommand) (BalanceResult, error) {
	card, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return BalanceResult{}, err
	}

	oldBalance := card.CurrentBalance()
	err = card.ApplyPayment(cmd.Amount, cmd.ReferenceId)

	if err != nil {
		return BalanceResult{}, err
	}

	err = h.store.Save(ctx, card)

	return BalanceResult{
		CardNo:     card.CardNo(),
		OldBalance: oldBalance,
		NewBalance: card.CurrentBalance(),
	}, err
}
