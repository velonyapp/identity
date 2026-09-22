package domain

import (
	"github.com/velonyapp/identity/internal/domain/service"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewUsernameAvailability,
	service.NewEmailAvailability,
)
