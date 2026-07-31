package journey

type PaymentPort interface {
	HasSufficientBalance(cardNo string, amount int) (bool, error)
	AuthorizePayment(cardNo string, referenceId string, amount int) error
}
