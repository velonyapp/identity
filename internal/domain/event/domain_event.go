package event

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	ID() uuid.UUID
	Type() string
	AggregateID() string
	OccurTime() time.Time
}

type BaseDomainEvent struct {
	id          uuid.UUID
	aggregateID string
	occurTime   time.Time
}

func NewBaseDomainEvent(aggregateID string, occurTime time.Time) BaseDomainEvent {
	return BaseDomainEvent{
		id:          uuid.Must(uuid.NewV7()),
		aggregateID: aggregateID,
		occurTime:   occurTime,
	}
}

func (e BaseDomainEvent) ID() uuid.UUID        { return e.id }
func (e BaseDomainEvent) AggregateID() string  { return e.aggregateID }
func (e BaseDomainEvent) OccurTime() time.Time { return e.occurTime }
