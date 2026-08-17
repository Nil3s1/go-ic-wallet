package domain

type PaymentPort interface {
	HasSufficientBalance(mediaId string, amount int) (bool, error)
	AuthorizePayment(mediaId string, referenceId string, amount int) error
}
