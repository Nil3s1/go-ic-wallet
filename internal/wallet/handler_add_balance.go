package wallet

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type AddBalanceHandler struct {
	store kernel.EventStore[*Card]
}

func NewAddBalanceHandler(store kernel.EventStore[*Card]) *AddBalanceHandler {
	return &AddBalanceHandler{
		store: store,
	}
}

func (h *AddBalanceHandler) Handle(ctx context.Context, cmd AddBalanceCommand) {

}
