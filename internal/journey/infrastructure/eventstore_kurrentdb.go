package infrastructure

import (
	"encoding/json"

	"github.com/Nil3s1/go-ic-wallet/internal/journey/domain"
	kernelApplication "github.com/Nil3s1/go-ic-wallet/internal/kernel/application"
	kernel "github.com/Nil3s1/go-ic-wallet/internal/kernel/domain"
	"github.com/Nil3s1/go-ic-wallet/internal/kernel/infrastructure/eventstore"
	"github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
)

var journeyLogEventTypes = eventstore.EventTypeRegistry{
	domain.EventJourneyLogCreatedV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyLogCreatedDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventJourneyStartedV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyStartedDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
	domain.EventJourneyEndedV1: func(data []byte) (kernel.DomainEvent, error) {
		var e domain.JourneyEndedDomainEventV1
		err := json.Unmarshal(data, &e)
		return e, err
	},
}

func NewKurrentDBStore(client *kurrentdb.Client) kernelApplication.EventStore[*domain.JourneyLog] {
	return eventstore.New(client, "journeyLog", domain.Rehydrate, journeyLogEventTypes)
}
