package journey

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type StartJourneyCommandHandler struct {
	store kernel.EventStore[*JourneyLog]
}

func NewStartJourneyCommandHandler(store kernel.EventStore[*JourneyLog]) *StartJourneyCommandHandler {
	return &StartJourneyCommandHandler{
		store: store,
	}
}

func (h *StartJourneyCommandHandler) Handle(ctx context.Context, cmd StartJourneyCommand) error {
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
