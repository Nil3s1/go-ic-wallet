package eventstore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

const readAllEventsCount = 10_000

// EventDecoder unmarshals a stored event's raw payload into its concrete DomainEvent value.
type EventDecoder func(data []byte) (kernel.DomainEvent, error)

// EventTypeRegistry maps a DomainEvent's EventName() to the decoder for its concrete type.
type EventTypeRegistry map[string]EventDecoder

// Store is a generic KurrentDB-backed kernel.EventStore[T] adapter, reused by every bounded context.
type Store[T kernel.HasDomainEvents] struct {
	client       *kurrentdb.Client
	streamPrefix string
	rehydrate    func([]kernel.DomainEvent) T
	eventTypes   EventTypeRegistry
}

// New builds a Store for T. rehydrate and eventTypes carry the bounded-context-specific knowledge.
func New[T kernel.HasDomainEvents](client *kurrentdb.Client, streamPrefix string, rehydrate func([]kernel.DomainEvent) T, eventTypes EventTypeRegistry) *Store[T] {
	return &Store[T]{
		client:       client,
		streamPrefix: streamPrefix,
		rehydrate:    rehydrate,
		eventTypes:   eventTypes,
	}
}

func (s *Store[T]) streamName(id string) string {
	return s.streamPrefix + "-" + id
}

func (s *Store[T]) Exists(ctx context.Context, id string) (bool, error) {
	stream, err := s.client.ReadStream(ctx, s.streamName(id), kurrentdb.ReadStreamOptions{
		Direction: kurrentdb.Forwards,
		From:      kurrentdb.Start{},
	}, 1)
	if err != nil {
		if isStreamNotFound(err) {
			return false, nil
		}
		return false, err
	}
	defer stream.Close()

	_, err = stream.Recv()
	if err != nil {
		if isStreamNotFound(err) || errors.Is(err, io.EOF) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Store[T]) Load(ctx context.Context, id string) (T, error) {
	var zero T

	stream, err := s.client.ReadStream(ctx, s.streamName(id), kurrentdb.ReadStreamOptions{
		Direction: kurrentdb.Forwards,
		From:      kurrentdb.Start{},
	}, readAllEventsCount)
	if err != nil {
		return zero, err
	}
	defer stream.Close()

	var events []kernel.DomainEvent
	for {
		resolved, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return zero, err
		}

		recorded := resolved.OriginalEvent()
		decode, ok := s.eventTypes[recorded.EventType]
		if !ok {
			return zero, fmt.Errorf("eventstore: unknown event type %q on stream %q", recorded.EventType, recorded.StreamID)
		}

		event, err := decode(recorded.Data)
		if err != nil {
			return zero, err
		}

		events = append(events, event)
	}

	if len(events) == 0 {
		return zero, fmt.Errorf("eventstore: stream %q not found", s.streamName(id))
	}

	return s.rehydrate(events), nil
}

func (s *Store[T]) Save(ctx context.Context, aggregate T) error {
	uncommitted := aggregate.UncommittedEvents()
	if len(uncommitted) == 0 {
		return nil
	}

	events := make([]kurrentdb.EventData, len(uncommitted))
	for i, e := range uncommitted {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}

		events[i] = kurrentdb.EventData{
			EventID:     uuid.New(),
			EventType:   e.EventName(),
			ContentType: kurrentdb.ContentTypeJson,
			Data:        data,
		}
	}

	// Version() only counts previously persisted events, so it doubles as the optimistic-concurrency check.
	var expectedState kurrentdb.StreamState
	if aggregate.Version() == 0 {
		expectedState = kurrentdb.NoStream{}
	} else {
		expectedState = kurrentdb.Revision(uint64(aggregate.Version() - 1))
	}

	_, err := s.client.AppendToStream(ctx, s.streamName(aggregate.Id()), kurrentdb.AppendToStreamOptions{
		StreamState: expectedState,
	}, events...)
	if err != nil {
		return err
	}

	aggregate.ClearUncommittedEvents()
	return nil
}

func isStreamNotFound(err error) bool {
	var esErr *kurrentdb.Error
	if errors.As(err, &esErr) {
		return esErr.Code() == kurrentdb.ErrorCodeResourceNotFound
	}
	return false
}
