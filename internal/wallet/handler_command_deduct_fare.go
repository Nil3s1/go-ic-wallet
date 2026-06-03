package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type DeductFareCommandHandler struct {
	store kernel.EventStore[*Card]
}

func NewDeductFareCommandHandler(store kernel.EventStore[*Card]) *DeductFareCommandHandler {
	return &DeductFareCommandHandler{
		store: store,
	}
}

func (h *DeductFareCommandHandler) Handle(ctx context.Context, cmd DeductFareCommand) error {
	card, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	err = card.DeductFare(cmd.Amount)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	return nil
}
