package application

import (
	"context"
)

type HasSufficientBalanceQueryHandler struct {
	repository CardProjectionRepository
}

func NewHasSufficientBalanceQueryHandler(repository CardProjectionRepository) *HasSufficientBalanceQueryHandler {
	return &HasSufficientBalanceQueryHandler{
		repository: repository,
	}
}

func (h *HasSufficientBalanceQueryHandler) Handle(ctx context.Context, query HasSufficientBalanceQuery) (bool, error) {
	return h.repository.HasSufficientBalance(ctx, query.CardNo, query.Amount)
}
