package application

import (
	"github.com/velonyapp/identity/internal/application/command"
	"github.com/velonyapp/identity/internal/application/domainevent"
	"github.com/velonyapp/identity/internal/application/query"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	command.NewBus,
	command.NewChangeUsernameHandler,
	command.NewRequestEmailChangeHandler,
	command.NewConfirmEmailChangeHandler,
	command.NewRemoveEmailHandler,
	command.NewRequestAvatarChangeHandler,
	command.NewConfirmAvatarChangeHandler,
	command.NewRemoveAvatarHandler,
	command.NewDeleteUserHandler,
	command.NewRegisterAuthHandler,
	command.NewLoginAuthHandler,
	command.NewRefreshAuthHandler,
	query.NewBus,
	query.NewGetUserHandler,
	query.NewBatchGetUsersHandler,
	domainevent.NewDispatcher,
	domainevent.NewUserCreatedHandler,
	domainevent.NewUserUpdatedHandler,
	domainevent.NewUserDeletedHandler,
)
