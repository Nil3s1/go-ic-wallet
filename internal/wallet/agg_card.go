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

	cardNo                string
	validTo               time.Time
	currentBalance        int //Currency in cents
	processedReferenceIds map[string]bool
}

func hasSufficientBalance(currentBalance int, amount int) bool {
	return currentBalance >= amount
}

func NewCard(initialBalance int) (*Card, error) {
	cardNo, err := generateCardNo()

	if err != nil {
		return nil, err
	}

	createdAt := time.Now().UTC()
	validTo := createdAt.AddDate(5, 0, 0)

	event := CardCreatedDomainEvent{
		CardNo:         cardNo,
		InitialBalance: initialBalance,
		CreatedAt:      createdAt,
		ValidTo:        validTo,
	}

	card := &Card{processedReferenceIds: make(map[string]bool)}
	card.ApplyEvent(event, card.applyEventFunction)

	return card, nil
}

func Rehydrate(events []kernel.DomainEvent) *Card {
	card := &Card{processedReferenceIds: make(map[string]bool)}
	card.BaseAggregate.LoadFromHistory(events, card.applyEventFunction)

	return card
}

func (c *Card) CardNo() string {
	return c.cardNo
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

	generatedReferenceID := uuid.New().String()

	event := BalanceAddedDomainEvent{
		BalanceAdded: int(value),
		ReferenceId:  generatedReferenceID,
	}

	c.ApplyEvent(event, c.applyEventFunction)

	return nil
}

func (c *Card) ApplyPayment(amount int, referenceID string) error {
	if !hasSufficientBalance(c.currentBalance, amount) {
		return errors.New("nicht genug Balance auf der Karte. Bitte Karte aufladen!")
	}

	event := ApplyPaymentDomainEvent{
		Amount:      int(amount),
		ReferenceId: referenceID,
	}

	c.ApplyEvent(event, c.applyEventFunction)

	return nil
}

func (c *Card) applyEventFunction(event kernel.DomainEvent) {
	switch e := event.(type) {
	case CardCreatedDomainEvent:
		c.SetId(e.CardNo)
		c.cardNo = e.CardNo
		c.SetCreatedAt(e.CreatedAt)
		c.validTo = e.ValidTo
		c.currentBalance = e.InitialBalance
	case BalanceAddedDomainEvent:
		c.currentBalance += e.BalanceAdded
	case ApplyPaymentDomainEvent:
		c.currentBalance -= e.Amount
		if e.ReferenceId != "" {
			c.processedReferenceIds[e.ReferenceId] = true
		}
	default:
	}
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
