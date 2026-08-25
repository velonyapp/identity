package transport

import (
	v1 "github.com/velony-app/identity/gen/api/v1"
	"github.com/velony-app/identity/internal/conf"
	"github.com/velony-app/identity/internal/presentation/api"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

func NewHTTPServer(
	c *conf.Transport,
	service *api.Service,
	tracingMiddleware TracesMiddleware,
	metricsMiddleware MetricsMiddleware,
	authMiddleware AuthMiddleware,
	validationMiddleware ValidationMiddleware,
) *http.Server {
	opts := []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			middleware.Middleware(tracingMiddleware),
			middleware.Middleware(metricsMiddleware),
			middleware.Middleware(authMiddleware),
			middleware.Middleware(validationMiddleware),
		),
	}

	if c.GetHttp().GetAddress() != "" {
		opts = append(opts, http.Address(c.GetHttp().GetAddress()))
	}
	if c.GetHttp().GetTimeout() != nil {
		opts = append(opts, http.Timeout(c.GetHttp().GetTimeout().AsDuration()))
	}

	srv := http.NewServer(opts...)

	v1.RegisterIdentityServiceHTTPServer(srv, service)

	return srv
}
