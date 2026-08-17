package application

import (
	"context"

	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type AddBalanceCommandHandler struct {
	store kernelApplication.EventStore[*domain.Card]
}

func NewAddBalanceHandler(store kernelApplication.EventStore[*domain.Card]) *AddBalanceCommandHandler {
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
