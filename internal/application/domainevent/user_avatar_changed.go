package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

var _ Handler[*event.UserAvatarChanged] = (*UserAvatarChangedHandler)(nil)

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
	var oldAvatarID, newAvatarID *string

	if domainEvent.OldAvatarID() != nil {
		value := domainEvent.OldAvatarID().String()
		oldAvatarID = &value
	}
	if domainEvent.NewAvatarID() != nil {
		value := domainEvent.NewAvatarID().String()
		newAvatarID = &value
	}

	integrationEvent := integrationevent.NewUserAvatarChanged(
		domainEvent.AggregateID(),
		oldAvatarID,
		newAvatarID,
		domainEvent.OccurTime(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
