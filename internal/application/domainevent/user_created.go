package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
)

type UserCreatedHandler struct {
	eventPublisher port.EventPublisher
}

func NewUserCreatedHandler(
	eventPublisher port.EventPublisher,
) *UserCreatedHandler {
	return &UserCreatedHandler{
		eventPublisher: eventPublisher,
	}
}

func (h *UserCreatedHandler) Execute(ctx context.Context, domainEvent event.UserCreated) error {
	var email, avatarKey *string

	if domainEvent.Email != nil {
		value := domainEvent.Email.Value()
		email = &value
	}
	if domainEvent.AvatarKey != nil {
		value := domainEvent.AvatarKey.String()
		avatarKey = &value
	}

	integrationEvent := integrationevent.NewUserCreated(
		domainEvent.AggregateID(),
		domainEvent.Username.Value(),
		email,
		domainEvent.FullName.Value(),
		avatarKey,
		domainEvent.CreateTime.Value(),
	)

	return h.eventPublisher.Publish(ctx, integrationEvent)
}
