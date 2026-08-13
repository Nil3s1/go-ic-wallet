package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type ApplyPaymentCommandHandler struct {
	store kernel.EventStore[*Card]
}

func NewApplyPaymentCommandHandler(store kernel.EventStore[*Card]) *ApplyPaymentCommandHandler {
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

	return err
}
