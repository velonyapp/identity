package presentation

import (
	"github.com/velony-app/identity/internal/presentation/api"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	api.NewService,
)
