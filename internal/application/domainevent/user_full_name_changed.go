package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"
)

type UserFullNameChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserFullNameChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserFullNameChangedHandler {
	return &UserFullNameChangedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserFullNameChangedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserFullNameChanged,
) error {
	payload := &v1.UserFullNameChangedPayload{
		OldFullName: domainEvent.OldFullName.Value(),
		NewFullName: domainEvent.NewFullName.Value(),
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		domainEvent.Type()+".v1",
		domainEvent.AggregateID(),
		domainEvent.AggregateType(),
		domainEvent.OccurTime(),
		payload,
	)
	if err != nil {
		return err
	}

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        integrationEvent,
		},
	)
}
