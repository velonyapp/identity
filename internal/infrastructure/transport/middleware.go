package transport

import (
	"context"
	"errors"

	v1 "github.com/velonyapp/identity/gen/api/v1"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/infrastructure/observability"

	kratosjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/selector"
	"github.com/go-kratos/kratos/v3/middleware/validate"
	"github.com/go-kratos/kratos/v3/transport"
	"github.com/golang-jwt/jwt/v5"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/protobuf/proto"
)

type AuthMiddleware middleware.Middleware

func NewAuthMiddleware(c *conf.Security) AuthMiddleware {
	auth := kratosjwt.Server(
		func(token *jwt.Token) (any, error) {
			return []byte(c.AccessToken.Secret), nil
		},
		kratosjwt.WithSigningMethod(jwt.SigningMethodHS256),
	)

	optionalAuth := optionalAuthentication(auth)

	return AuthMiddleware(
		middleware.Chain(
			selector.Server(optionalAuth).
				Path(
					"/velony.identity.api.v1.IdentityService/GetUser",
				).
				Build(),
			selector.Server(auth).
				Path(
					"/velony.identity.api.v1.IdentityService/UpdateUser",
					"/velony.identity.api.v1.IdentityService/RequestUserEmailChange",
					"/velony.identity.api.v1.IdentityService/RequestUserAvatarChange",
					"/velony.identity.api.v1.IdentityService/DeleteUser",
				).
				Build(),
		),
	)
}

func optionalAuthentication(auth middleware.Middleware) middleware.Middleware {
	return func(next middleware.Handler) middleware.Handler {
		authenticated := auth(next)

		return func(ctx context.Context, req any) (any, error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return next(ctx, req)
			}

			if tr.RequestHeader().Get("Authorization") == "" {
				return next(ctx, req)
			}

			return authenticated(ctx, req)
		}
	}
}

type MetricsMiddleware middleware.Middleware

func NewMetricsMiddleware(serverMetrics *observability.ServerMetrics) MetricsMiddleware {
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
	)
}
