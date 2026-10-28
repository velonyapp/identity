package infrastructure

import (
	"github.com/velonyapp/identity/internal/infrastructure/auth"
	"github.com/velonyapp/identity/internal/infrastructure/data/mysql"
	"github.com/velonyapp/identity/internal/infrastructure/gateway/asset"
	"github.com/velonyapp/identity/internal/infrastructure/messaging/protobuf"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewUserRepo,
	mysql.NewSessionRepo,
	mysql.NewAvatarRepo,
	mysql.NewUserQuery,
	mysql.NewEventPublisher,
	mysql.NewJobPublisher,
	protobuf.NewEventEncoder,
	protobuf.NewJobEncoder,
	auth.NewPasswordManager,
	auth.NewEmailChangeRequestTokenManager,
	auth.NewAccessTokenManager,
	auth.NewSessionTokenManager,
	asset.NewClient,
	asset.NewService,
)
