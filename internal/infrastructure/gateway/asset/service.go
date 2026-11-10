package asset

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/entity"
	"github.com/velonyapp/identity/internal/domain/vo"

	"go.einride.tech/aip/resourcename"
)

const (
	imageResourcePattern = "images/{image}"
)

var _ port.AssetService = (*service)(nil)

type service struct {
	client v1.AssetServiceClient
	log    *slog.Logger
}

func NewService(
	client v1.AssetServiceClient,
	log *slog.Logger,
) port.AssetService {
	return &service{client: client, log: log}
}

func (s *service) CreateAvatar(ctx context.Context, userID vo.UserID) (*entity.Avatar, string, error) {
	var suffix [3]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		s.log.ErrorContext(ctx,
			"failed to generate avatar object key suffix",
			slog.Any("error", err),
		)
		return nil, "", err
	}

	result, err := s.client.CreateImage(ctx, &v1.CreateImageRequest{
		Image: &v1.Image{
			Tags:      []string{"identity", "user"},
			Metadata:  map[string]string{"user_id": userID.String()},
			ObjectKey: fmt.Sprintf("users/%s/avatar-%s.webp", userID.String(), hex.EncodeToString(suffix[:])),
			WidthPx:   512,
			HeightPx:  512,
			Format:    v1.Image_WEBP,
		},
	})
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to create avatar image",
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
		return nil, "", err
	}

	imageID, err := parseImageName(result.Name)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to parse avatar image resource name",
			slog.String("user_id", userID.String()),
			slog.String("image_name", result.Name),
			slog.Any("error", err),
		)
		return nil, "", err
	}

	avatarID, err := vo.NewAvatarID(imageID)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to create avatar ID value object",
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
		return nil, "", err
	}
	key, err := vo.NewAvatarKey(result.ObjectKey)
	if err != nil {
		s.log.ErrorContext(ctx,
			"failed to create avatar key value object",
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
		return nil, "", err
	}

	avatar := entity.NewAvatar(
		avatarID,
		userID,
		key,
	)
	uploadURL := *result.UploadUrl

	return avatar, uploadURL, nil
}

func parseImageName(name string) (string, error) {
	var imageID string
	if err := resourcename.Sscan(name, imageResourcePattern, &imageID); err != nil {
		return "", err
	}

	return imageID, nil
}
