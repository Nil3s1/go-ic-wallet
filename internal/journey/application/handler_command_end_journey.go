package application

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/domain"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
)

type EndJourneyCommandHandler struct {
	store       application.EventStore[*domain.JourneyLog]
	paymentPort domain.PaymentPort
	calculator  domain.FareCalculator
}

func NewEndJourneyCommandHandler(paymentPort domain.PaymentPort, store application.EventStore[*domain.JourneyLog], calculator domain.FareCalculator) *EndJourneyCommandHandler {
	return &EndJourneyCommandHandler{
		store:       store,
		paymentPort: paymentPort,
		calculator:  calculator,
	}
}

func (h *EndJourneyCommandHandler) Handle(ctx context.Context, cmd EndJourneyCommand) error {
	jl, err := h.store.Load(ctx, cmd.MediaId)
	if err != nil {
		return err
	}

	cf, err := h.calculator.CalculateFare(jl.StartStation(), cmd.EndStation)
	if err != nil {
		return err
	}

	err = h.paymentPort.AuthorizePayment(ctx, cmd.MediaId, jl.JourneyReferenceId(), cf.Fare())
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
