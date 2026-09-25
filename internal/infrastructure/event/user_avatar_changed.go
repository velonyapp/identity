package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userAvatarKeyChangedPayload(event integrationevent.UserAvatarChanged) *eventv1.UserAvatarChangedPayload {
	return &eventv1.UserAvatarChangedPayload{
		OldAvatarKey: event.OldAvatarKey(),
		NewAvatarKey: event.NewAvatarKey(),
	}
}
