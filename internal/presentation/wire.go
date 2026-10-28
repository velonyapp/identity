package presentation

import (
	"github.com/velonyapp/identity/internal/presentation/middleware"
	"github.com/velonyapp/identity/internal/presentation/service"
	"github.com/velonyapp/identity/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAPI,
	service.NewConsumer,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewKafkaConsumer,
	middleware.NewTracing,
	middleware.NewMetrics,
	middleware.NewErrorMapper,
	middleware.NewAuthentication,
	middleware.NewValidation,
)
