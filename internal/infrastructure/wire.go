package infrastructure

import (
	"github.com/velony-app/identity/internal/infrastructure/auth"
	"github.com/velony-app/identity/internal/infrastructure/data/mysql"
	"github.com/velony-app/identity/internal/infrastructure/data/redis"
	"github.com/velony-app/identity/internal/infrastructure/observability"
	"github.com/velony-app/identity/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewUserRepo,
	mysql.NewSessionRepo,
	mysql.NewOutboxPublisher,
	redis.NewConnection,
	redis.NewCache,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewTracesMiddleware,
	transport.NewMetricsMiddleware,
	transport.NewAuthMiddleware,
	transport.NewValidationMiddleware,
	observability.NewMetrics,
	auth.NewPasswordHasher,
	auth.NewTokenProvider,
)
