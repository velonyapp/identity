package command

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

func NewBus(
	changeUsernameHandler ChangeUsernameHandler,
	requestEmailChangeHandler RequestEmailChangeHandler,
	confirmEmailChangeHandler ConfirmEmailChangeHandler,
	removeEmailHandler RemoveEmailHandler,
	requestAvatarChangeHandler RequestAvatarChangeHandler,
	confirmAvatarChangeHandler ConfirmAvatarChangeHandler,
	removeAvatarHandler RemoveAvatarHandler,
	deleteUserHandler DeleteUserHandler,
	registerAuthHandler RegisterAuthHandler,
	loginAuthHandler LoginAuthHandler,
	refreshAuthHandler RefreshAuthHandler,
) *Bus {
	bus := &Bus{}

	Register(bus, changeUsernameHandler)
	Register(bus, requestEmailChangeHandler)
	Register(bus, confirmEmailChangeHandler)
	Register(bus, removeEmailHandler)
	Register(bus, requestAvatarChangeHandler)
	Register(bus, confirmAvatarChangeHandler)
	Register(bus, removeAvatarHandler)
	Register(bus, deleteUserHandler)
	Register(bus, registerAuthHandler)
	Register(bus, loginAuthHandler)
	Register(bus, refreshAuthHandler)

	return bus
}

type Bus struct {
	handlers sync.Map
}

func Register[R any, C Command[R]](bus *Bus, handler Handler[C, R]) {
	t := typeOf[C]()

	bus.handlers.Store(t, erasedHandler(func(ctx context.Context, cmd any) (any, error) {
		return handler.Handle(ctx, cmd.(C))
	}))
}

func Send[R any, C Command[R]](ctx context.Context, bus *Bus, cmd C) (R, error) {
	var zero R

	raw, ok := bus.handlers.Load(typeOf[C]())
	if !ok {
		return zero, fmt.Errorf("no handler registered for %T", cmd)
	}

	result, err := raw.(erasedHandler)(ctx, cmd)
	if err != nil {
		return zero, err
	}

	return result.(R), nil
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

type erasedHandler func(context.Context, any) (any, error)
