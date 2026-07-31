package wallet

type HasSufficientBalanceQueryHandler struct {
	repository CardProjectionRepository
}

func NewHasSufficientBalanceQueryHandler(repository CardProjectionRepository) *HasSufficientBalanceQueryHandler {
	return &HasSufficientBalanceQueryHandler{
		repository: repository,
	}
}

func (h *HasSufficientBalanceQueryHandler) Handle(query HasSufficientBalanceQuery) (bool, error) {
	return h.repository.HasSufficientBalance(query.CardNo, query.Amount)
}
