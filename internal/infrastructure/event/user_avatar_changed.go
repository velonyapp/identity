package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userAvatarChangedPayload(event integrationevent.UserAvatarChanged) *eventv1.UserAvatarChangedPayload {
	return &eventv1.UserAvatarChangedPayload{
		OldAvatarId: event.OldAvatarID(),
		NewAvatarId: event.NewAvatarID(),
	}
}
