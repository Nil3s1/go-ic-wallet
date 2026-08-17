package application

import (
	"context"
	"errors"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/domain"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type StartJourneyCommandHandler struct {
	store       kernel.EventStore[*domain.JourneyLog]
	paymentPort domain.PaymentPort
}

func NewStartJourneyCommandHandler(paymentPort domain.PaymentPort, store kernel.EventStore[*domain.JourneyLog]) *StartJourneyCommandHandler {
	return &StartJourneyCommandHandler{
		store:       store,
		paymentPort: paymentPort,
	}
}

func (h *StartJourneyCommandHandler) Handle(ctx context.Context, cmd StartJourneyCommand) error {
	result, err := h.paymentPort.HasSufficientBalance(cmd.MediaId, domain.BaseFare)

	if err != nil {
		return err
	}

	if result == false {
		return errors.New("Insufficient Balance to Start the Journey")
	}

	exists, err := h.store.Exists(ctx, cmd.MediaId)

	if err != nil {
		return err
	}

	var jl *domain.JourneyLog

	if exists {
		jl, err = h.store.Load(ctx, cmd.MediaId)
	} else {
		jl, err = domain.NewJourneyLog(cmd.MediaId)
	}

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
