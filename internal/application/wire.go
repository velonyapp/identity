package application

import (
	"github.com/velonyapp/identity/internal/application/command"
	"github.com/velonyapp/identity/internal/application/domainevent"
	"github.com/velonyapp/identity/internal/application/query"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	query.NewGetUserHandler,
	query.NewBatchGetUsersHandler,
	command.NewRegisterAuthHandler,
	command.NewLoginAuthHandler,
	command.NewRefreshAuthHandler,
	command.NewUpdateUserHandler,
	command.NewDeleteUserHandler,
	domainevent.NewDispatcher,
	domainevent.NewUserCreatedHandler,
	domainevent.NewUserUsernameChangedHandler,
	domainevent.NewUserFullNameChangedHandler,
	domainevent.NewUserAvatarKeyChangedHandler,
	domainevent.NewUserDeletedHandler,
)
