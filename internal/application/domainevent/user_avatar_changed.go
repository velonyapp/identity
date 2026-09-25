package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserAvatarChangedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserAvatarChangedHandler(
	eventPublisher port.EventPublisher,
) *UserAvatarChangedHandler {
	return &UserAvatarChangedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserAvatarChangedHandler) Execute(ctx context.Context, domainEvent *event.UserAvatarChanged) error {
	var oldAvatarKey, newAvatarKey *string

	if domainEvent.OldAvatarKey() != nil {
		value := domainEvent.OldAvatarKey().String()
		oldAvatarKey = &value
	}
	if domainEvent.NewAvatarKey() != nil {
		value := domainEvent.NewAvatarKey().String()
		newAvatarKey = &value
	}

	integrationEvent := integrationevent.NewUserAvatarChanged(
		domainEvent.AggregateID(),
		oldAvatarKey,
		newAvatarKey,
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
