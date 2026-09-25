package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"

)

func sessionRevokedPayload(event integrationevent.SessionRevoked) *eventv1.SessionRevokedPayload {
	return &eventv1.SessionRevokedPayload{
	}
}
