//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"log/slog"

	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/application"
	"github.com/velonyapp/identity/internal/domain"
	"github.com/velonyapp/identity/internal/info"
	"github.com/velonyapp/identity/internal/infrastructure"
	"github.com/velonyapp/identity/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	context.Context,
	*info.Service,
	*conf.Data,
	*conf.Transport,
	*conf.Auth,
	*conf.Observability,
	*conf.Gateway,
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
