package application

import (
	"context"

	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
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
	card, err := h.repository.GetByKey(ctx, query.CardNo)
	if err != nil {
		return false, err
	}

	return domain.HasSufficientBalance(card.CurrentBalance, query.Amount), nil
}
