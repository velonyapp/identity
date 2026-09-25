package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userCreatedPayload(event integrationevent.UserCreated) *eventv1.UserCreatedPayload {
	return &eventv1.UserCreatedPayload{
		Username: event.Username,
		FullName: event.FullName,
	}
}
