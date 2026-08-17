package application

type CreateCardCommand struct {
	CardNo         string
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
