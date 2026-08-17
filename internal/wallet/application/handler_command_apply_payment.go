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

func (h *ApplyPaymentCommandHandler) Handle(ctx context.Context, cmd ApplyPaymentCommand) error {
	card, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	err = card.ApplyPayment(cmd.Amount, cmd.ReferenceId)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	return nil
}
