package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type CreateCardCommandHandler struct {
	store kernel.EventStore[*Card]
}

func NewCreateCardHandler(store kernel.EventStore[*Card]) *CreateCardCommandHandler {
	return &CreateCardCommandHandler{
		store: store,
	}
}

func (h *CreateCardCommandHandler) Handle(ctx context.Context, cmd CreateCardCommand) error {
	card, err := NewCard(cmd.InitialBalance)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	if err != nil {
		return err
	}

	return err
}
