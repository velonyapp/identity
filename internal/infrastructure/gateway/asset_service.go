package gateway

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	assetv1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/identity/internal/application/port"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AssetService struct {
	client assetv1.AssetServiceClient
}

func NewAssetService(client assetv1.AssetServiceClient) port.AssetService {
	return &AssetService{client: client}
}

func (s *AssetService) PresignAvatar(ctx context.Context, userID string) (string, error) {
	var tokenBytes [8]byte

	if _, err := rand.Read(tokenBytes[:]); err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes[:])

	width := uint32(512)
	height := uint32(512)
	quality := uint32(85)

	result, err := s.client.PresignImage(ctx,
		&assetv1.PresignImageRequest{
			StorageKey: fmt.Sprintf(
				"users/%s/avatar-%s.webp",
				userID,
				token,
			),
			Transform: &assetv1.ImageTransform{
				Resize: &assetv1.ImageResize{
					Width:        &width,
					Height:       &height,
					Fit:          assetv1.ImageResizeFit_IMAGE_RESIZE_FIT_COVER,
					Gravity:      assetv1.ImageGravity_IMAGE_GRAVITY_CENTER,
					AllowUpscale: true,
				},
				Encoding: &assetv1.ImageEncoding{
					Format:  assetv1.ImageFormat_IMAGE_FORMAT_WEBP,
					Quality: &quality,
				},
			},
			ExpireTime: timestamppb.New(
				time.Now().Add(10 * time.Minute),
			),
		},
	)
	if err != nil {
		return "", err
	}

	return result.UploadUrl, nil
}
