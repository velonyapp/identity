package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userUsernameChangedPayload(event integrationevent.UserUsernameChanged) *eventv1.UserUsernameChangedPayload {
	return &eventv1.UserUsernameChangedPayload{
		OldUsername: event.OldUsername(),
		NewUsername: event.NewUsername(),
	}
}
