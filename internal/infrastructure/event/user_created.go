package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userCreatedPayload(event integrationevent.UserCreated) *eventv1.UserCreatedPayload {
	return &eventv1.UserCreatedPayload{
		Username:   event.Username,
		Email:      event.Email,
		FullName:   event.FullName,
		AvatarKey:  event.AvatarKey,
		CreateTime: timestamppb.New(event.CreateTime),
	}
}
