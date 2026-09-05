package domainevent

import (
	"context"

	v1 "github.com/velonyapp/identity/gen/event/v1"
	integrationevent "github.com/velonyapp/identity/internal/application/event"
	"github.com/velonyapp/identity/internal/application/port"
	domainevent "github.com/velonyapp/identity/internal/domain/event"

	"google.golang.org/protobuf/types/known/timestamppb"
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
	payload := &v1.UserAvatarKeyChangedPayload{
		UserId:     domainEvent.AggregateID(),
		AvatarKey:  domainEvent.AvatarKey.Value(),
		UpdateTime: timestamppb.New(domainEvent.UpdateTime.Value()),
	}

	integrationEvent, err := integrationevent.NewIntegrationEvent(
		"user.full-name.changed.v1",
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
