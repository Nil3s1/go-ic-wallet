package wallet

type CardModelRepository interface {
	GetCard(cardNo string) (CardModel, error)

	Update(card CardModel)
	Add(card CardModel)
}
