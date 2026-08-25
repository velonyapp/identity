package port

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	GetMany(ctx context.Context, keys []string, dest map[string]any) (map[string]bool, error)

	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	SetMany(ctx context.Context, items map[string]any, ttl time.Duration) error

	Delete(ctx context.Context, key string) error
	DeleteMany(ctx context.Context, keys []string) error
}
