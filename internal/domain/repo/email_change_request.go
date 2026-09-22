package repo

import (
	"context"

	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"
)

type EmailChangeRequest interface {
	FindByToken(ctx context.Context, token vo.RequestToken) (*entity.EmailChangeRequest, error)

	Save(ctx context.Context, request *entity.EmailChangeRequest) error
}
