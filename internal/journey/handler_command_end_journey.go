package journey

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type EndJourneyCommandHandler struct {
	store kernel.EventStore[*JourneyLog]
}

func NewEndJourneyCommandHandler(store kernel.EventStore[*JourneyLog]) *EndJourneyCommandHandler {
	return &EndJourneyCommandHandler{
		store: store,
	}
}

func (h *EndJourneyCommandHandler) Handle(ctx context.Context, cmd EndJourneyCommand) error {
	jl, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	err = jl.StartJourney(cmd.EndStation)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, jl)

	return err
}
