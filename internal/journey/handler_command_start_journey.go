package journey

import (
	"context"
	"errors"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type StartJourneyCommandHandler struct {
	store       kernel.EventStore[*JourneyLog]
	paymentPort PaymentPort
}

func NewStartJourneyCommandHandler(paymentPort PaymentPort, store kernel.EventStore[*JourneyLog]) *StartJourneyCommandHandler {
	return &StartJourneyCommandHandler{
		store:       store,
		paymentPort: paymentPort,
	}
}

func (h *StartJourneyCommandHandler) Handle(ctx context.Context, cmd StartJourneyCommand) error {
	result, err := h.paymentPort.HasSufficientBalance(cmd.CardNo, BaseFare)

	if err != nil {
		return err
	}

	if result == false {
		return errors.New("Insufficient Balance to Start the Journey")
	}

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
