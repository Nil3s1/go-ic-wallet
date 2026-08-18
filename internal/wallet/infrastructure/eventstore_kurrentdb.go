package infrastructure

import (
	"encoding/json"

	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	kernel "github.com/Nil3s1/go-ic-wallet/internal/kernel/domain"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/eventstore"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var cardEventTypes = eventstore.EventTypeRegistry{
	domain.EventCardCreatedV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.CardCreatedDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventBalanceAddedV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.BalanceAddedDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventApplyPaymentV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.ApplyPaymentDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
}

func NewKurrentDBStore(client *kurrentdb.Client) kernelApplication.EventStore[*domain.Card] {
	return eventstore.New(client, "card", domain.Rehydrate, cardEventTypes)
}
