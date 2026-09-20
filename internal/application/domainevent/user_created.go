package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/application/integrationevent"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/event"
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

	return h.outboxPublisher.PublishMessage(ctx, integrationEvent)
}
