package application

import "github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"

type HasSufficientBalanceQueryHandler struct {
	repository domain.CardProjectionRepository
}

func NewHasSufficientBalanceQueryHandler(repository domain.CardProjectionRepository) *HasSufficientBalanceQueryHandler {
	return &HasSufficientBalanceQueryHandler{
		repository: repository,
	}
}

func (h *HasSufficientBalanceQueryHandler) Handle(query HasSufficientBalanceQuery) (bool, error) {
	return h.repository.HasSufficientBalance(query.CardNo, query.Amount)
}
