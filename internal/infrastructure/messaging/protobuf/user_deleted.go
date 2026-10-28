package protobuf

import (
	v1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userDeletedPayload(event integrationevent.UserDeleted) *v1.UserDeletedPayload {
	return &v1.UserDeletedPayload{
		UserId:  event.AggregateID(),
		Version: event.Version(),
	}
}
