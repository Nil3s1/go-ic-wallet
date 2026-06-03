package journey

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type CreateJourneyLogCommandHandler struct {
	store kernel.EventStore[*JourneyLog]
}

func NewCreateJourneyLogCommandHandler(store kernel.EventStore[*JourneyLog]) *CreateJourneyLogCommandHandler {
	return &CreateJourneyLogCommandHandler{
		store: store,
	}
}

func (h *CreateJourneyLogCommandHandler) Handle(ctx context.Context, cmd CreateJourneyLogCommand) error {
	var err error
	var jl *JourneyLog

	if cmd.PredecessorCardNo == "" {
		jl, err = NewJourneyLog(cmd.CardNo)
	} else {
		jl, err = NewJourneyLogWithPredecessor(cmd.CardNo, cmd.PredecessorCardNo)
	}

	if err != nil {
		return err
	}

	err = h.store.Save(ctx, jl)

	return err
}
