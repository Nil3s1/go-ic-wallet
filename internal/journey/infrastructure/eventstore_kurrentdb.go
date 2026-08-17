package infrastructure

import (
	"encoding/json"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/domain"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/eventstore"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var journeyLogEventTypes = eventstore.EventTypeRegistry{
	domain.EventJourneyLogCreated: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyLogCreatedDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventJourneyStarted: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyStartedDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventJourneyEnded: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyEndedDomainEvent
		err := json.Unmarshal(data, &e)
		return e, err
	},
}

func NewKurrentDBStore(client *kurrentdb.Client) kernel.EventStore[*domain.JourneyLog] {
	return eventstore.New(client, "journeyLog", domain.Rehydrate, journeyLogEventTypes)
}
