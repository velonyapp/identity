package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/velonyapp/identity/internal/application/port"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	client *redis.Client
}

type noopCache struct{}

func NewCache(client *redis.Client) port.Cache {
	if client == nil {
		return &noopCache{}
	}

	return &cache{client: client}
}

func (r *cache) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}

	return true, nil
}

func (r *cache) GetMany(ctx context.Context, keys []string, dest map[string]any) (map[string]bool, error) {
	found := make(map[string]bool, len(keys))

	if len(keys) == 0 {
		return found, nil
	}

	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	for i, value := range values {
		key := keys[i]

		if value == nil {
			found[key] = false
			continue
		}

		target, ok := dest[key]
		if !ok {
			found[key] = true
			continue
		}

		data, ok := value.(string)
		if !ok {
			continue
		}

		if err := json.Unmarshal([]byte(data), target); err != nil {
			return nil, err
		}

		found[key] = true
	}

	return found, nil
}

func (r *cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *cache) SetMany(ctx context.Context, items map[string]any, ttl time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	values := make([]any, 0, len(items)*2)

	for key, value := range items {
		data, err := json.Marshal(value)
		if err != nil {
			return err
		}

		values = append(values, key, data)
	}

	args := redis.MSetEXArgs{}

	if ttl > 0 {
		args.Expiration = &redis.ExpirationOption{
			Mode:  redis.PX,
			Value: ttl.Milliseconds(),
		}
	}

	return r.client.MSetEX(ctx, args, values...).Err()
}

func (r *cache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *cache) DeleteMany(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(ctx, keys...).Err()
}

func (n *noopCache) Get(ctx context.Context, key string, dest any) (bool, error) {
	return false, nil
}

func (n *noopCache) GetMany(ctx context.Context, keys []string, dest map[string]any) (map[string]bool, error) {
	return make(map[string]bool, len(keys)), nil
}

func (n *noopCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return nil
}

func (n *noopCache) SetMany(ctx context.Context, items map[string]any, ttl time.Duration) error {
	return nil
}

func (n *noopCache) Delete(ctx context.Context, key string) error {
	return nil
}

func (n *noopCache) DeleteMany(ctx context.Context, keys []string) error {
	return nil
}
