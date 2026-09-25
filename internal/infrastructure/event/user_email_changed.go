package event

import (
	eventv1 "github.com/velonyapp/identity/gen/event/v1"
	"github.com/velonyapp/identity/internal/application/integrationevent"
)

func userEmailChangedPayload(event integrationevent.UserEmailChanged) *eventv1.UserEmailChangedPayload {
	return &eventv1.UserEmailChangedPayload{
		OldEmail:   event.OldEmail(),
		NewEmail:   event.NewEmail(),
	}
}
