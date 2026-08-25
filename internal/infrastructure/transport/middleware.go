package transport

import (
	"github.com/velony-app/identity/internal/conf"
	"github.com/velony-app/identity/internal/infrastructure/observability"

	kratosjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/selector"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/golang-jwt/jwt/v5"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

type AuthMiddleware middleware.Middleware

func NewAuthMiddleware(c *conf.Auth) AuthMiddleware {
	auth := kratosjwt.Server(
		func(token *jwt.Token) (any, error) {
			return []byte(c.GetAccessToken().GetSecret()), nil
		},
		kratosjwt.WithSigningMethod(jwt.SigningMethodHS256),
	)

	return AuthMiddleware(
		selector.Server(auth).
			Path(
				"/velony.identity.api.v1.IdentityService/UpdateUser",
				"/velony.identity.api.v1.IdentityService/ReplaceUserAvatar",
				"/velony.identity.api.v1.IdentityService/DeleteUser",
			).
			Build(),
	)
}

type MetricsMiddleware middleware.Middleware

func NewMetricsMiddleware(serverMetrics *observability.Metrics) MetricsMiddleware {
	return MetricsMiddleware(
		metrics.Server(
			metrics.WithSeconds(serverMetrics.Seconds),
			metrics.WithRequests(serverMetrics.Requests),
		),
	)
}

type TracesMiddleware middleware.Middleware

func NewTracesMiddleware() TracesMiddleware {
	return TracesMiddleware(tracing.Server())
}

type ValidationMiddleware middleware.Middleware

func NewValidationMiddleware() ValidationMiddleware {
	return ValidationMiddleware(
		validate.Validator(func(req any) error {
			message, ok := req.(proto.Message)
			if !ok {
				return nil
			}

			return fieldbehavior.ValidateRequiredFields(message)
		}),
	)
}
