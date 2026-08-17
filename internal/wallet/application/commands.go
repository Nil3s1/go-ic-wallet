package application

type CreateCardCommand struct {
	InitialBalance uint
}

type AddBalanceCommand struct {
	CardNo string
	Amount uint
}

type ApplyPaymentCommand struct {
	CardNo      string
	Amount      uint
	ReferenceId string
}
