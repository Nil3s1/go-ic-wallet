package kurrentdbstore

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Nil3s1/go-ic-wallet/internal/journey"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet"
)

type EventCodec interface {
	Encode(event kernel.DomainEvent) (string, []byte, error)
	Decode(eventType string, data []byte) (kernel.DomainEvent, error)
}

type JSONCodec struct {
	types map[string]reflect.Type
}

func NewJSONCodec() *JSONCodec {
	return &JSONCodec{types: make(map[string]reflect.Type)}
}

func (c *JSONCodec) Register(events ...kernel.DomainEvent) {
	for _, event := range events {
		if event == nil {
			continue
		}

		c.types[event.EventName()] = reflect.TypeOf(event)
	}
}

func (c *JSONCodec) Encode(event kernel.DomainEvent) (string, []byte, error) {
	if event == nil {
		return "", nil, fmt.Errorf("cannot encode nil event")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return "", nil, err
	}

	return event.EventName(), data, nil
}

func (c *JSONCodec) Decode(eventType string, data []byte) (kernel.DomainEvent, error) {
	t, ok := c.types[eventType]
	if !ok {
		return nil, fmt.Errorf("event type %q is not registered", eventType)
	}

	value := reflect.New(t)
	if err := json.Unmarshal(data, value.Interface()); err != nil {
		return nil, err
	}

	event, ok := value.Elem().Interface().(kernel.DomainEvent)
	if !ok {
		return nil, fmt.Errorf("event type %q does not implement kernel.DomainEvent", eventType)
	}

	return event, nil
}

func NewWalletEventCodec() *JSONCodec {
	codec := NewJSONCodec()
	codec.Register(
		wallet.CardCreatedDomainEvent{},
		wallet.BalanceAddedDomainEvent{},
		wallet.ApplyPaymentDomainEvent{},
	)

	return codec
}

func NewJourneyEventCodec() *JSONCodec {
	codec := NewJSONCodec()
	codec.Register(
		journey.JourneyLogCreatedDomainEvent{},
		journey.JourneyStartedDomainEvent{},
		journey.JourneyEndedDomainEvent{},
	)

	return codec
}
