package redis

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/velonyapp/identity/internal/conf"
)

func NewConnection(c *conf.Data) (*redis.Client, error) {
	if c.Redis == nil {
		return nil, nil
	}

	cfg, err := redis.ParseURL(c.Redis.Dsn)
	if err != nil {
		return nil, err
	}

	if c.Redis.MaxActiveConnections != nil {
		cfg.MaxActiveConns = int(*c.Redis.MaxActiveConnections)
	}

	client := redis.NewClient(cfg)

	if err := redisotel.InstrumentTracing(client); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}
