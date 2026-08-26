//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/velony-app/identity/internal/application"
	"github.com/velony-app/identity/internal/conf"
	"github.com/velony-app/identity/internal/domain"
	"github.com/velony-app/identity/internal/infrastructure"
	"github.com/velony-app/identity/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Data,
	*conf.Transport,
	*conf.Auth,
	*conf.Observability,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		presentation.ProviderSet,
		infrastructure.ProviderSet,
		application.ProviderSet,
		domain.ProviderSet,
		newApp,
	))
}
