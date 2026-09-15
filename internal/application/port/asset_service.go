package port

import "context"

type AssetService interface {
	PresignAvatar(ctx context.Context, userID string) (string, error)
}
