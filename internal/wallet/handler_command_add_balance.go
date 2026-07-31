package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type AddBalanceCommandHandler struct {
	store kernel.EventStore[*Card]
}

func NewAddBalanceHandler(store kernel.EventStore[*Card]) *AddBalanceCommandHandler {
	return &AddBalanceCommandHandler{
		store: store,
	}
}

func (h *AddBalanceCommandHandler) Handle(ctx context.Context, cmd AddBalanceCommand) error {
	card, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	err = card.AddBalance(cmd.Amount)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	return err
}
