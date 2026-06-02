package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type CreateCardHandler struct {
	store kernel.EventStore[*Card]
}

func NewCreateCardHandler(store kernel.EventStore[*Card]) *CreateCardHandler {
	return &CreateCardHandler{
		store: store,
	}
}

func (h *CreateCardHandler) Handle(ctx context.Context, cmd CreateCardCommand) error {
	card, err := NewCard(cmd.InitialBalance)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, card)

	return err
}
