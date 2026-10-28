package domainevent

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/event"
)

type Handler[T event.DomainEvent] interface {
	Execute(ctx context.Context, event T) error
}
