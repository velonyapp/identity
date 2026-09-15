package infrastructure

import (
	"github.com/velonyapp/identity/internal/infrastructure/auth"
	"github.com/velonyapp/identity/internal/infrastructure/data/mysql"
	"github.com/velonyapp/identity/internal/infrastructure/data/redis"
	"github.com/velonyapp/identity/internal/infrastructure/gateway"
	"github.com/velonyapp/identity/internal/infrastructure/observability"
	"github.com/velonyapp/identity/internal/infrastructure/transport"

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
	observability.NewServerMetrics,
	observability.NewOpenTelemetry,
	auth.NewPasswordHasher,
	auth.NewTokenProvider,
	gateway.NewAssetClient,
	gateway.NewAssetService,
)
