//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"log/slog"

	"github.com/velonyapp/identity/internal/application"
	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/infrastructure"
	"github.com/velonyapp/identity/internal/presentation"

	"github.com/go-kratos/kratos/v3"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(
	*conf.Auth,
	*conf.Transport,
	*conf.Data,
	*conf.Gateway,
	*slog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		application.ProviderSet,
		infrastructure.ProviderSet,
		presentation.ProviderSet,
		newApp,
	))
}
