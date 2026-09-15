package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/port"
)

type PresignUserAvatar struct {
	UserID string
}

type PresignUserAvatarResult struct {
	UploadURL string
}

type PresignUserAvatarHandler struct {
	assetService port.AssetService
}

func NewPresignUserAvatarHandler(
	assetService port.AssetService,
) *PresignUserAvatarHandler {
	return &PresignUserAvatarHandler{
		assetService: assetService,
	}
}

func (h *PresignUserAvatarHandler) Execute(
	ctx context.Context,
	cmd *PresignUserAvatar,
) (*PresignUserAvatarResult, error) {
	uploadURL, err := h.assetService.PresignAvatar(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &PresignUserAvatarResult{
		UploadURL: uploadURL,
	}, nil
}
