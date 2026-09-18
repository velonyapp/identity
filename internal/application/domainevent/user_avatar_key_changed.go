package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
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
	domainEvent event.UserAvatarKeyChanged,
) error {
	var oldAvatarKey, newAvatarKey *string

	if domainEvent.OldAvatarKey != nil {
		value := domainEvent.OldAvatarKey.String()
		oldAvatarKey = &value
	}
	if domainEvent.NewAvatarKey != nil {
		value := domainEvent.NewAvatarKey.String()
		newAvatarKey = &value
	}

	integrationEvent := integrationevent.NewUserAvatarKeyChanged(
		domainEvent.AggregateID(),
		oldAvatarKey,
		newAvatarKey,
		domainEvent.UpdateTime.Value(),
	)

	return h.outboxPublisher.PublishMessage(ctx,
		port.OutboxMessage{
			PartitionKey: domainEvent.AggregateID(),
			Event:        integrationEvent,
		},
	)
}
