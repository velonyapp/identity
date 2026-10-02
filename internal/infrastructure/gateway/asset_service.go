package gateway

import (
	"context"
	"fmt"
	"time"

	assetv1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/identity/internal/application/port"
	"github.com/velonyapp/identity/internal/domain/vo"
	"go.einride.tech/aip/resourcename"

	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	avatarResourcePattern = "images/{image}"
)

var _ port.AssetService = (*assetService)(nil)

type assetService struct {
	client assetv1.AssetServiceClient
}

func NewAssetService(client assetv1.AssetServiceClient) port.AssetService {
	return &assetService{client: client}
}

func (s *assetService) CreateAvatar(ctx context.Context, userID vo.UserID) (vo.AvatarID, error) {
	result, err := s.client.CreateImage(ctx, &assetv1.CreateImageRequest{
		Image: &assetv1.Image{
			Tags:      []string{"identity", "user"},
			ObjectKey: fmt.Sprintf("users/%s/avatar-%s.webp", userID.Value(), "test"),
		},
	})
	if err != nil {
		return vo.AvatarID{}, err
	}

	var avatarID string
	if err := resourcename.Sscan(result.Name, avatarResourcePattern, &avatarID); err != nil {
		return vo.AvatarID{}, err
	}

	return vo.NewAvatarID(avatarID), nil
}

func (s *assetService) PresignAvatar(ctx context.Context, avatarID vo.AvatarID, ttl time.Duration) (string, error) {
	result, err := s.client.PresignImage(ctx, &assetv1.PresignImageRequest{
		Name: resourcename.Sprint(avatarResourcePattern, avatarID.Value()),
		Ttl:  durationpb.New(ttl),
	})
	if err != nil {
		return "", err
	}

	return result.UploadUrl, nil
}
