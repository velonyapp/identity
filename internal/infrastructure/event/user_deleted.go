package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userDeletedPayload(event integrationevent.UserDeleted) *eventv1.UserDeletedPayload {
	return &eventv1.UserDeletedPayload{}
}
