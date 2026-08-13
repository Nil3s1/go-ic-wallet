package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
)

type stubCardStore struct {
	card    *Card
	saveErr error
}

func (s *stubCardStore) Exists(_ context.Context, _ string) (bool, error) {
	return true, nil
}

func (s *stubCardStore) Load(_ context.Context, _ string) (*Card, error) {
	return s.card, nil
}

func (s *stubCardStore) Save(_ context.Context, _ *Card) error {
	return s.saveErr
}

func TestApplyPaymentCommandHandlerReturnsSaveError(t *testing.T) {
	card, err := NewCard(200)
	if err != nil {
		t.Fatalf("NewCard() error = %v", err)
	}
	card.ClearUncommittedEvents()

	expectedErr := errors.New("save failed")
	handler := NewApplyPaymentCommandHandler(&stubCardStore{
		card:    card,
		saveErr: expectedErr,
	})

	err = handler.Handle(context.Background(), ApplyPaymentCommand{
		CardNo:      card.CardNo(),
		Amount:      100,
		ReferenceId: "ref-1",
	})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Handle() error = %v, want %v", err, expectedErr)
	}
}

var _ kernel.EventStore[*Card] = (*stubCardStore)(nil)
