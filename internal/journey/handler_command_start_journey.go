package journey

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type StartJourneyCommandHandler struct {
	store    kernel.EventStore[*JourneyLog]
	provider CardProvider
}

func NewStartJourneyCommandHandler(provider CardProvider, store kernel.EventStore[*JourneyLog]) *StartJourneyCommandHandler {
	return &StartJourneyCommandHandler{
		store:    store,
		provider: provider,
	}
}

func (h *StartJourneyCommandHandler) Handle(ctx context.Context, cmd StartJourneyCommand) error {
	CardHasSufficientBalance(h.provider, cmd.CardNo, BaseFare)
	jl, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	err = jl.StartJourney(cmd.StartStation)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, jl)

	return err
}
