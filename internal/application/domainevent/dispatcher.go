package domainevent

import (
	"context"
	"fmt"
	"reflect"

	"github.com/velonyapp/identity/internal/domain/event"
)

type Dispatcher struct {
	handlers map[reflect.Type][]func(
		context.Context,
		event.DomainEvent,
	) error
}

func NewDispatcher(
	userCreatedHandler *UserCreatedHandler,
	userUsernameChangedHandler *UserUsernameChangedHandler,
	userEmailChangedHandler *UserEmailChangedHandler,
	userFullNameChangedHandler *UserFullNameChangedHandler,
	userAvatarKeyChangedHandler *UserAvatarKeyChangedHandler,
	userDeletedHandler *UserDeletedHandler,
	sessionCreatedHandler *SessionCreatedHandler,
	sessionRefreshedHandler *SessionRefreshedHandler,
	sessionRevokedHandler *SessionRevokedHandler,
) *Dispatcher {
	dispatcher := &Dispatcher{
		handlers: make(
			map[reflect.Type][]func(
				context.Context,
				event.DomainEvent,
			) error,
		),
	}

	registerHandler(dispatcher, userCreatedHandler)
	registerHandler(dispatcher, userUsernameChangedHandler)
	registerHandler(dispatcher, userEmailChangedHandler)
	registerHandler(dispatcher, userFullNameChangedHandler)
	registerHandler(dispatcher, userAvatarKeyChangedHandler)
	registerHandler(dispatcher, userDeletedHandler)
	registerHandler(dispatcher, sessionCreatedHandler)
	registerHandler(dispatcher, sessionRefreshedHandler)
	registerHandler(dispatcher, sessionRevokedHandler)

	return dispatcher
}

func (d *Dispatcher) Dispatch(ctx context.Context, event event.DomainEvent) error {
	handlers := d.handlers[reflect.TypeOf(event)]

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}

	return nil
}

func registerHandler[T event.DomainEvent](d *Dispatcher, h Handler[T]) {
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	d.handlers[eventType] = append(
		d.handlers[eventType],
		func(ctx context.Context, event event.DomainEvent) error {
			typedEvent, ok := event.(T)
			if !ok {
				return fmt.Errorf(
					"invalid domain event type: expected %v, got %T",
					eventType,
					event,
				)
			}

			return h.Execute(ctx, typedEvent)
		},
	)
}
