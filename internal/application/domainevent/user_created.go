package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"
)

type UserCreatedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserCreatedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserCreatedHandler {
	return &UserCreatedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserCreatedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserCreated,
) error {
	payload := &v1.UserCreatedPayload{
		Username: domainEvent.Username.Value(),
		FullName: domainEvent.FullName.Value(),
	}

	if domainEvent.Email != nil {
		value := domainEvent.Email.Value()
		payload.Email = &value
	}
	if domainEvent.AvatarKey != nil {
		value := domainEvent.AvatarKey.String()
		payload.AvatarKey = &value
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
