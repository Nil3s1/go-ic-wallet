package kurrentdbstore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/Nil3s1/go-ic-wallet/internal/journey"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet"
	"github.com/google/uuid"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

const (
	WalletStreamPrefix  = "wallet-card"
	JourneyStreamPrefix = "journey-log"
)

type streamReader interface {
	Recv() (*kurrentdb.ResolvedEvent, error)
	Close()
}

type streamClient interface {
	AppendToStream(context.Context, string, kurrentdb.AppendToStreamOptions, ...kurrentdb.EventData) (*kurrentdb.WriteResult, error)
	ReadStream(context.Context, string, kurrentdb.ReadStreamOptions, uint64) (streamReader, error)
}

type kurrentClient struct {
	client *kurrentdb.Client
}

func (c kurrentClient) AppendToStream(
	ctx context.Context,
	streamID string,
	opts kurrentdb.AppendToStreamOptions,
	events ...kurrentdb.EventData,
) (*kurrentdb.WriteResult, error) {
	return c.client.AppendToStream(ctx, streamID, opts, events...)
}

func (c kurrentClient) ReadStream(
	ctx context.Context,
	streamID string,
	opts kurrentdb.ReadStreamOptions,
	count uint64,
) (streamReader, error) {
	return c.client.ReadStream(ctx, streamID, opts, count)
}

type Store[T kernel.HasDomainEvents] struct {
	client       streamClient
	streamPrefix string
	codec        EventCodec
	rehydrate    func([]kernel.DomainEvent) T
}

func NewStore[T kernel.HasDomainEvents](
	client *kurrentdb.Client,
	streamPrefix string,
	codec EventCodec,
	rehydrate func([]kernel.DomainEvent) T,
) *Store[T] {
	return newStoreWithClient[T](kurrentClient{client: client}, streamPrefix, codec, rehydrate)
}

func NewWalletStore(client *kurrentdb.Client) *Store[*wallet.Card] {
	return NewStore(client, WalletStreamPrefix, NewWalletEventCodec(), wallet.Rehydrate)
}

func NewJourneyStore(client *kurrentdb.Client) *Store[*journey.JourneyLog] {
	return NewStore(client, JourneyStreamPrefix, NewJourneyEventCodec(), journey.Rehydrate)
}

func newStoreWithClient[T kernel.HasDomainEvents](
	client streamClient,
	streamPrefix string,
	codec EventCodec,
	rehydrate func([]kernel.DomainEvent) T,
) *Store[T] {
	return &Store[T]{
		client:       client,
		streamPrefix: streamPrefix,
		codec:        codec,
		rehydrate:    rehydrate,
	}
}

func (s *Store[T]) Exists(ctx context.Context, id string) (bool, error) {
	reader, err := s.client.ReadStream(ctx, s.streamID(id), readOptions(), 1)
	if err != nil {
		if isNotFound(err) {
			return false, nil
		}

		return false, err
	}
	defer reader.Close()

	_, err = reader.Recv()
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, io.EOF):
		return false, nil
	case isNotFound(err):
		return false, nil
	default:
		return false, err
	}
}

func (s *Store[T]) Load(ctx context.Context, id string) (T, error) {
	var zero T

	reader, err := s.client.ReadStream(ctx, s.streamID(id), readOptions(), math.MaxUint64)
	if err != nil {
		return zero, err
	}
	defer reader.Close()

	events := make([]kernel.DomainEvent, 0)

	for {
		resolved, err := reader.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return zero, err
		}
		if resolved == nil || resolved.Event == nil {
			continue
		}

		event, err := s.codec.Decode(resolved.Event.EventType, resolved.Event.Data)
		if err != nil {
			return zero, err
		}

		events = append(events, event)
	}

	if len(events) == 0 {
		return zero, fmt.Errorf("aggregate %q has no events", id)
	}

	return s.rehydrate(events), nil
}

func (s *Store[T]) Save(ctx context.Context, aggregate T) error {
	events := aggregate.UncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	appendEvents := make([]kurrentdb.EventData, 0, len(events))
	for _, event := range events {
		eventType, data, err := s.codec.Encode(event)
		if err != nil {
			return err
		}

		appendEvents = append(appendEvents, kurrentdb.EventData{
			EventID:     uuid.New(),
			EventType:   eventType,
			ContentType: kurrentdb.ContentTypeJson,
			Data:        data,
		})
	}

	_, err := s.client.AppendToStream(
		ctx,
		s.streamID(aggregate.Id()),
		kurrentdb.AppendToStreamOptions{StreamState: expectedStreamState(aggregate.Version())},
		appendEvents...,
	)
	if err != nil {
		return err
	}

	for range events {
		aggregate.IncrementVersion()
	}
	aggregate.ClearUncommittedEvents()

	return nil
}

func expectedStreamState(version int) kurrentdb.StreamState {
	if version == 0 {
		return kurrentdb.NoStream{}
	}

	return kurrentdb.Revision(uint64(version - 1))
}

func readOptions() kurrentdb.ReadStreamOptions {
	return kurrentdb.ReadStreamOptions{
		Direction: kurrentdb.Forwards,
		From:      kurrentdb.Start{},
	}
}

func (s *Store[T]) streamID(id string) string {
	if s.streamPrefix == "" {
		return id
	}

	return fmt.Sprintf("%s-%s", s.streamPrefix, id)
}

func isNotFound(err error) bool {
	var esErr *kurrentdb.Error
	if !errors.As(err, &esErr) {
		return false
	}

	return esErr.IsErrorCode(kurrentdb.ErrorCodeResourceNotFound)
}
