package wallet

type CreateCardCommand struct {
	CardNo         string
	InitialBalance int
}

type AddBalanceCommand struct {
	CardNo string
	Amount int
}

type ApplyPaymentCommand struct {
	CardNo      string
	Amount      int
	ReferenceId string
}
