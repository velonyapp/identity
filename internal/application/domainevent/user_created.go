package domainevent

import (
	"context"

	v1 "github.com/velony-app/identity/gen/event/v1"
	integrationevent "github.com/velony-app/identity/internal/application/event"
	"github.com/velony-app/identity/internal/application/port"
	domainevent "github.com/velony-app/identity/internal/domain/event"

	"google.golang.org/protobuf/types/known/timestamppb"
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
		UserId:     domainEvent.AggregateID(),
		Username:   domainEvent.Username.Value(),
		FullName:   domainEvent.FullName.Value(),
		CreateTime: timestamppb.New(domainEvent.CreateTime.Value()),
	}
	if domainEvent.AvatarKey != nil {
		value := domainEvent.AvatarKey.Value()
		payload.AvatarKey = &value
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		"user.created.v1",
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
