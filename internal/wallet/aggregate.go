package wallet

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/google/uuid"
)

type Card struct {
	kernel.BaseAggregate

	cardNo         string
	validTo        time.Time
	currentBalance int //Currency in cents

}

func NewCard() (*Card, error) {
	cardNo, err := generateCardNo()

	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()
	validTo := createdAt.AddDate(5, 0, 0)

	event := CardCreatedDomainEvent{
		Id:        uuid.NewString(),
		CardNo:    cardNo,
		CreatedAt: createdAt,
		ValidTo:   validTo,
	}

	card := &Card{}
	card.applyEvent(event)

	return card, nil
}

func Rehydrate(events []kernel.DomainEvent) *Card {
	card := &Card{}
	for _, e := range events {
		card.applyEvent(e)
	}

	card.ClearUncommittedEvents()
	return card
}

func (c *Card) ValidTo() time.Time {
	return c.validTo
}

func (c *Card) CurrentBalance() int {
	return c.currentBalance
}

func (c *Card) AddBalance(value int) error {
	if value <= 0 {
		return errors.New("Betrag muss größer als 0 sein")
	}

	event := BalanceAddedDomainEvent{
		BalanceAdded: int(value),
	}

	c.applyEvent(event)

	return nil
}

func (c *Card) DeductFare(value int) error {
	if c.currentBalance < value {
		return errors.New("nicht genug Blance auf der Karte. Bitte Karte aufladen!")
	}
	event := FareDeductedDomainEvent{
		DeductedFare: int(value),
	}

	c.applyEvent(event)

	return nil
}

func (c *Card) applyEvent(event kernel.DomainEvent) {
	switch e := event.(type) {
	case CardCreatedDomainEvent:
		c.SetId(e.Id)
		c.cardNo = e.CardNo
		c.SetCreatedAt(e.CreatedAt)
		c.validTo = e.ValidTo
		c.currentBalance = 0
	case BalanceAddedDomainEvent:
		c.currentBalance += e.BalanceAdded
	case FareDeductedDomainEvent:
		c.currentBalance -= e.DeductedFare
	default:
	}

	c.AddEvent(event)
}

func generateCardNo() (string, error) {
	min := big.NewInt(10000000000)

	rangeLimit := big.NewInt(90000000000)

	n, err := rand.Int(rand.Reader, rangeLimit)
	if err != nil {
		return "", err
	}

	cardNumber := new(big.Int).Add(min, n)

	return cardNumber.String(), nil
}
