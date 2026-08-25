package event

import "github.com/google/uuid"

type DomainEvent interface {
	ID() uuid.UUID
	Type() string
	AggregateID() string
}

type BaseDomainEvent struct {
	id          uuid.UUID
	aggregateID string
}

func NewBaseDomainEvent(aggregateID string) BaseDomainEvent {
	return BaseDomainEvent{
		id:          uuid.Must(uuid.NewV7()),
		aggregateID: aggregateID,
	}
}

func (e BaseDomainEvent) ID() uuid.UUID       { return e.id }
func (e BaseDomainEvent) AggregateID() string { return e.aggregateID }
