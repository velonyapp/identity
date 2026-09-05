package presentation

import (
	"github.com/velonyapp/identity/internal/presentation/api"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
