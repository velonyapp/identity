package transport

import (
	"errors"

	v1 "github.com/velonyapp/identity/gen/api/v1"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/presentation/api"
	"github.com/velonyapp/identity/internal/presentation/observability"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

func NewGRPCServer(
	c *conf.Transport,
	service *api.Service,
	serverMetrics *observability.ServerMetrics,
) *grpc.Server {
	opts := []grpc.ServerOption{
		grpc.Address(c.Grpc.Address),
		grpc.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metrics.Server(
				metrics.WithSeconds(serverMetrics.Seconds),
				metrics.WithRequests(serverMetrics.Requests),
			),
			validate.Validator(func(req any) error {
				switch req := req.(type) {
				case *v1.UpdateUserRequest:
					if req.GetUser() == nil {
						return errors.New("missing required field: user")
					}

					return fieldbehavior.ValidateRequiredFieldsWithMask(
						req.GetUser(),
						req.GetUpdateMask(),
					)

				default:
					message, ok := req.(proto.Message)
					if !ok {
						return nil
					}

					return fieldbehavior.ValidateRequiredFields(message)
				}
			}),
		),
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)

	v1.RegisterIdentityServiceServer(srv, service)

	return srv
}
