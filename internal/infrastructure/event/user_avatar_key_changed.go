package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func userAvatarKeyChangedPayload(event integrationevent.UserAvatarKeyChanged) *eventv1.UserAvatarKeyChangedPayload {
	return &eventv1.UserAvatarKeyChangedPayload{
		OldAvatarKey: event.OldAvatarKey,
		NewAvatarKey: event.NewAvatarKey,
		UpdateTime:   timestamppb.New(event.UpdateTime),
	}
}
