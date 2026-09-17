package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"
)

type UserAvatarKeyChangedHandler struct {
	outboxPublisher port.OutboxPublisher
}

func NewUserAvatarKeyChangedHandler(
	outboxPublisher port.OutboxPublisher,
) *UserAvatarKeyChangedHandler {
	return &UserAvatarKeyChangedHandler{
		outboxPublisher: outboxPublisher,
	}
}

func (h *UserAvatarKeyChangedHandler) Execute(
	ctx context.Context,
	domainEvent domainevent.UserAvatarKeyChanged,
) error {
	var oldAvatarKey *string
	if domainEvent.OldAvatarKey != nil {
		value := domainEvent.OldAvatarKey.String()
		oldAvatarKey = &value
	}

	var newAvatarKey *string
	if domainEvent.NewAvatarKey != nil {
		value := domainEvent.NewAvatarKey.String()
		newAvatarKey = &value
	}

	payload := &v1.UserAvatarKeyChangedPayload{
		OldAvatarKey: oldAvatarKey,
		NewAvatarKey: newAvatarKey,
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
