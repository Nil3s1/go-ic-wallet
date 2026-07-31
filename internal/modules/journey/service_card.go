package journey

type CardProvider interface {
	GetCurrentBalance(cardNo string) (int, error)
}

func CardHasSufficientBalance(provider CardProvider, cardNo string, amount int) (bool, error) {
	balance, err := provider.GetCurrentBalance(cardNo)

	if err != nil {
		return false, err
	}

	return balance >= amount, nil
}
