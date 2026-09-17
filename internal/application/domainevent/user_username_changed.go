package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"
)

type UserUsernameChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserUsernameChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserUsernameChangedHandler {
	return &UserUsernameChangedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserUsernameChangedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserUsernameChanged,
) error {
	payload := &v1.UserUsernameChangedPayload{
		OldUsername: domainEvent.OldUsername.Value(),
		NewUsername: domainEvent.NewUsername.Value(),
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
