package wallet

type CreateCardCommand struct {
	CardNo         string
	InitialBalance int
}

type AddBalanceCommand struct {
	CardNo string
	Amount int
}

type DeductFareCommand struct {
	CardNo string
	Amount int
}
