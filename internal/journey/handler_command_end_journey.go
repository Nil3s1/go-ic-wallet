package journey

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type EndJourneyCommandHandler struct {
	store      kernel.EventStore[*JourneyLog]
	calculator FareCalculator
}

func NewEndJourneyCommandHandler(store kernel.EventStore[*JourneyLog], calculator FareCalculator) *EndJourneyCommandHandler {
	return &EndJourneyCommandHandler{
		store:      store,
		calculator: calculator,
	}
}

func (h *EndJourneyCommandHandler) Handle(ctx context.Context, cmd EndJourneyCommand) error {
	jl, err := h.store.Load(ctx, cmd.CardNo)

	if err != nil {
		return err
	}

	cf, err := h.calculator.CalculateFare(jl.StartStation(), cmd.EndStation)

	if err != nil {
		return err
	}

	err = jl.EndJourney(cmd.EndStation, cf)

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, jl)

	return err
}
