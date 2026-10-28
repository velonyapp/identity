package transport

import (
	v1 "github.com/velonyapp/identity/gen/api/v1"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/presentation/middleware"
	"github.com/velonyapp/identity/internal/presentation/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func NewHTTPServer(
	c *conf.Transport,
	service *service.API,
	tracing middleware.Tracing,
	metrics middleware.Metrics,
	errorMapper middleware.ErrorMapper,
	authentication middleware.Authentication,
	validation middleware.Validation,
) *http.Server {
	opts := []http.ServerOption{
		http.Address(c.Http.Address),
		http.Middleware(
			recovery.Recovery(),
			tracing.Middleware(),
			metrics.Middleware(),
			errorMapper.Middleware(),
			authentication.Middleware(),
			validation.Middleware(),
		),
	}

	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterIdentityServiceHTTPServer(srv, service)

	return srv
}
