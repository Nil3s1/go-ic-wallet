package wallet

type BalanceAdapter struct {
	repo CardModelRepository
}

func NewBalanceAdapter(repo CardModelRepository) *BalanceAdapter {
	return &BalanceAdapter{
		repo: repo,
	}
}

func (ba *BalanceAdapter) CardHasSufficientBalance(cardNo string, amount int) (bool, error) {
	card, err := ba.repo.GetCard(cardNo)

	if err != nil {
		return false, err
	}

	return hasSufficientBalance(card.CurrentBalance, amount), nil
}
