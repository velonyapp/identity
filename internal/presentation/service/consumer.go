package service

import (
	"context"

	"github.com/velonyapp/identity/internal/application/command"
)

type Consumer struct {
	commandBus *command.Bus
}

func NewConsumer(
	commandBus *command.Bus,
) *Consumer {
	return &Consumer{
		commandBus: commandBus,
	}
}

type ConfirmUserAvatarChangeRequest struct {
	UserID string
}

func (s *Consumer) ConfirmUserAvatarChange(ctx context.Context, req ConfirmUserAvatarChangeRequest) error {
	if _, err := command.Send(ctx, s.commandBus, command.ConfirmAvatarChange{
		UserID: req.UserID,
	}); err != nil {
		return err
	}

	return nil
}
