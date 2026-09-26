package command

import (
	"context"

	"github.com/velonyapp/identity/internal/application/port"
)

type RequestUserAvatarChange struct {
	UserID string
}

type RequestUserAvatarChangeResult struct {
	UploadURL string
}

type RequestUserAvatarChangeHandler struct {
	assetService port.AssetService
}

func NewRequestUserAvatarChangeHandler(
	assetService port.AssetService,
) *RequestUserAvatarChangeHandler {
	return &RequestUserAvatarChangeHandler{
		assetService: assetService,
	}
}

func (h *RequestUserAvatarChangeHandler) Execute(
	ctx context.Context,
	cmd *RequestUserAvatarChange,
) (*RequestUserAvatarChangeResult, error) {
	uploadURL, err := h.assetService.PresignAvatar(ctx, cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &RequestUserAvatarChangeResult{
		UploadURL: uploadURL,
	}, nil
}
