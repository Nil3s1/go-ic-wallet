package kurrentdbstore

import (
	"context"
	"io"
	"testing"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

type fakeStreamClient struct {
	streams         map[string][]kurrentdb.EventData
	appendRevisions []kurrentdb.StreamState
}

func newFakeStreamClient() *fakeStreamClient {
	return &fakeStreamClient{streams: make(map[string][]kurrentdb.EventData)}
}

func (c *fakeStreamClient) AppendToStream(
	_ context.Context,
	streamID string,
	opts kurrentdb.AppendToStreamOptions,
	events ...kurrentdb.EventData,
) (*kurrentdb.WriteResult, error) {
	c.appendRevisions = append(c.appendRevisions, opts.StreamState)
	c.streams[streamID] = append(c.streams[streamID], events...)

	return &kurrentdb.WriteResult{}, nil
}

func (c *fakeStreamClient) ReadStream(
	_ context.Context,
	streamID string,
	_ kurrentdb.ReadStreamOptions,
	_ uint64,
) (streamReader, error) {
	return &fakeReadStream{events: c.streams[streamID]}, nil
}

type fakeReadStream struct {
	events []kurrentdb.EventData
	index  int
}

func (s *fakeReadStream) Recv() (*kurrentdb.ResolvedEvent, error) {
	if s.index >= len(s.events) {
		return nil, io.EOF
	}

	event := s.events[s.index]
	s.index++

	return &kurrentdb.ResolvedEvent{
		Event: &kurrentdb.RecordedEvent{
			EventType: event.EventType,
			Data:      event.Data,
		},
	}, nil
}

func (s *fakeReadStream) Close() {}

func TestWalletStoreSavesAndLoadsEvents(t *testing.T) {
	client := newFakeStreamClient()
	store := newStoreWithClient[*wallet.Card](client, WalletStreamPrefix, NewWalletEventCodec(), wallet.Rehydrate)

	card, err := wallet.NewCard(100)
	if err != nil {
		t.Fatalf("NewCard() error = %v", err)
	}

	if err := store.Save(context.Background(), card); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, ok := client.appendRevisions[0].(kurrentdb.NoStream); !ok {
		t.Fatalf("first append should expect no stream, got %T", client.appendRevisions[0])
	}

	loaded, err := store.Load(context.Background(), card.CardNo())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.CurrentBalance() != 100 {
		t.Fatalf("loaded balance = %d, want 100", loaded.CurrentBalance())
	}

	if err := loaded.AddBalance(50); err != nil {
		t.Fatalf("AddBalance() error = %v", err)
	}

	if err := store.Save(context.Background(), loaded); err != nil {
		t.Fatalf("second Save() error = %v", err)
	}

	revision, ok := client.appendRevisions[1].(kurrentdb.StreamRevision)
	if !ok {
		t.Fatalf("second append should expect stream revision, got %T", client.appendRevisions[1])
	}
	if revision.Value != 0 {
		t.Fatalf("second append revision = %d, want 0", revision.Value)
	}

	reloaded, err := store.Load(context.Background(), card.CardNo())
	if err != nil {
		t.Fatalf("second Load() error = %v", err)
	}

	if reloaded.CurrentBalance() != 150 {
		t.Fatalf("reloaded balance = %d, want 150", reloaded.CurrentBalance())
	}
}

func TestWalletStoreExistsReturnsFalseForMissingStream(t *testing.T) {
	client := newFakeStreamClient()
	store := newStoreWithClient[*wallet.Card](client, WalletStreamPrefix, NewWalletEventCodec(), wallet.Rehydrate)

	exists, err := store.Exists(context.Background(), "missing")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Fatal("Exists() = true, want false")
	}
}

type testEvent struct {
	Value string
}

func (e testEvent) EventName() string {
	return "test-event"
}

func TestJSONCodecRejectsUnknownEventTypes(t *testing.T) {
	codec := NewJSONCodec()
	codec.Register(testEvent{})

	_, err := codec.Decode("other-event", []byte(`{"value":"x"}`))
	if err == nil {
		t.Fatal("Decode() error = nil, want error")
	}
}

var _ kernel.DomainEvent = testEvent{}
