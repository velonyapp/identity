package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func sessionCreatedPayload(event integrationevent.SessionCreated) *eventv1.SessionCreatedPayload {
	return &eventv1.SessionCreatedPayload{
		UserId: event.UserID(),
	}
}
