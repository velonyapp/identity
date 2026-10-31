package middleware

import (
	"context"

	"github.com/velonyapp/identity/internal/conf"

	kjwt "github.com/go-kratos/kratos/contrib/middleware/jwt/v3"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/selector"
	"github.com/golang-jwt/jwt/v5"
)

type Authentication middleware.Middleware

func NewAuthentication(c *conf.Auth) (Authentication, error) {
	protected := map[string]struct{}{
		"/velony.identity.api.v1.IdentityService/DeleteUser":              {},
		"/velony.identity.api.v1.IdentityService/ChangeUserUsername":      {},
		"/velony.identity.api.v1.IdentityService/RequestUserEmailChange":  {},
		"/velony.identity.api.v1.IdentityService/ConfirmUserEmailChange":  {},
		"/velony.identity.api.v1.IdentityService/RemoveUserEmail":         {},
		"/velony.identity.api.v1.IdentityService/RequestUserAvatarChange": {},
		"/velony.identity.api.v1.IdentityService/RemoveUserAvatar":        {},
	}
	optional := map[string]struct{}{
		"/velony.identity.api.v1.IdentityService/GetUser": {},
	}

	auth := kjwt.Server(
		func(_ *jwt.Token) (any, error) {
			return []byte(c.AccessToken.Secret), nil
		},
		kjwt.WithParserOptions(
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		),
	)
	optionalAuth := func(next middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			called := false

			reply, err := auth(func(ctx context.Context, req any) (any, error) {
				called = true
				return next(ctx, req)
			})(ctx, req)

			if !called && err != nil {
				return next(ctx, req)
			}

			return reply, err
		}
	}

	selected := middleware.Chain(
		selector.Server(auth).
			Match(func(_ context.Context, operation string) bool {
				_, ok := protected[operation]
				return ok
			}).
			Build(),
		selector.Server(optionalAuth).
			Match(func(_ context.Context, operation string) bool {
				_, ok := optional[operation]
				return ok
			}).
			Build(),
	)

	return Authentication(selected), nil
}

func (m Authentication) Middleware() middleware.Middleware {
	return middleware.Middleware(m)
}
