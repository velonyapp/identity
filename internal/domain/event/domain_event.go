package event

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	ID() string
	AggregateID() string
	OccurTime() time.Time
}

type BaseDomainEvent struct {
	id          string
	aggregateID string
	occurTime   time.Time
}

func NewBaseDomainEvent(aggregateID string, occurTime time.Time) BaseDomainEvent {
	return BaseDomainEvent{
		id:          uuid.Must(uuid.NewV7()).String(),
		aggregateID: aggregateID,
		occurTime:   occurTime,
	}
}

func (e *BaseDomainEvent) ID() string           { return e.id }
func (e *BaseDomainEvent) AggregateID() string  { return e.aggregateID }
func (e *BaseDomainEvent) OccurTime() time.Time { return e.occurTime }
