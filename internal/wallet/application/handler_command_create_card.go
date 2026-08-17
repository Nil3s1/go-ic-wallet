package application

import (
	"context"

	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
)

type CreateCardCommandHandler struct {
	store kernelApplication.EventStore[*domain.Card]
}

func NewCreateCardHandler(store kernelApplication.EventStore[*domain.Card]) *CreateCardCommandHandler {
	return &CreateCardCommandHandler{
		store: store,
	}
}

func (h *CreateCardCommandHandler) Handle(ctx context.Context, cmd CreateCardCommand) (CreateCardResult, error) {
	card, err := domain.NewCard(cmd.InitialBalance)

	if err != nil {
		return CreateCardResult{}, err
	}

	err = h.store.Save(ctx, card)

	if err != nil {
		return CreateCardResult{}, err
	}

	return CreateCardResult{
		CardNo:  card.CardNo(),
		ValidTo: card.ValidTo(),
	}, nil
}
