package infrastructure

import (
	"github.com/velonyapp/identity/internal/infrastructure/data/mysql"
	"github.com/velonyapp/identity/internal/infrastructure/data/redis"
	"github.com/velonyapp/identity/internal/infrastructure/event"
	"github.com/velonyapp/identity/internal/infrastructure/gateway"
	"github.com/velonyapp/identity/internal/infrastructure/observability"
	"github.com/velonyapp/identity/internal/infrastructure/security"
	"github.com/velonyapp/identity/internal/infrastructure/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	event.NewEncoder,
	mysql.NewConnection,
	mysql.NewUnitOfWork,
	mysql.NewUserRepo,
	mysql.NewSessionRepo,
	mysql.NewEventPublisher,
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
	security.NewPasswordHasher,
	security.NewAuthTokenProvider,
	security.NewRequestVerifier,
	gateway.NewAssetClient,
	gateway.NewAssetService,
)
