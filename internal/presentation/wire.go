package presentation

import (
	"github.com/velonyapp/identity/internal/presentation/api"
	"github.com/velonyapp/identity/internal/presentation/observability"
	"github.com/velonyapp/identity/internal/presentation/transport"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
	transport.NewGRPCServer,
	transport.NewHTTPServer,
	transport.NewKafkaConsumer,
	observability.NewServerMetrics,
)
