package infrastructure

import (
	"encoding/json"

	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/eventstore"
	"github.com/Nil3s1/go-ic-wallet/internal/wallet/domain"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var cardEventTypes = eventstore.EventTypeRegistry{
	domain.EventCardCreated: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.CardCreatedDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventBalanceAdded: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.BalanceAddedDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventApplyPayment: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.ApplyPaymentDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
}

func NewKurrentDBStore(client *kurrentdb.Client) kernel.EventStore[*domain.Card] {
	return eventstore.New(client, "card", domain.Rehydrate, cardEventTypes)
}
