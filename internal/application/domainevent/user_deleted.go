package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserDeletedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserDeletedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserDeletedHandler {
	return &UserDeletedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserDeletedHandler) Execute(ctx context.Context, domainEvent event.UserDeleted) error {
	integrationEvent := integrationevent.NewUserDeleted(
		domainEvent.AggregateID(),
		domainEvent.DeleteTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        integrationEvent,
		},
	)
}
