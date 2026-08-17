package domain

type CardProjectionRepository interface {
	GetCard(cardNo string) (CardProjection, error)
	HasSufficientBalance(cardNo string, amount int) (bool, error)

	Update(card CardProjection) error
	Add(card CardProjection) error
}
