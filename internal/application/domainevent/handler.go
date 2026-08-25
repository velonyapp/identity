package domainevent

import (
	"context"

	"github.com/velony-app/identity/internal/domain/event"
)

type Handler[T event.DomainEvent] interface {
	Execute(ctx context.Context, event T) error
}
