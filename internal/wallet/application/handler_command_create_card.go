package application

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type CreateCardCommandHandler struct {
	store kernel.EventStore[*domain.Card]
}

func NewCreateCardHandler(store kernel.EventStore[*domain.Card]) *CreateCardCommandHandler {
	return &CreateCardCommandHandler{
		store: store,
	}
}

func (h *CreateCardCommandHandler) Handle(ctx context.Context, cmd CreateCardCommand) error {
	card, err := domain.NewCard(cmd.InitialBalance)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	if err != nil {
		return err
	}

	return err
}
