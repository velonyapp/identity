package redis

import (
	"github.com/velony-app/identity/internal/conf"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func NewConnection(c *conf.Data) (*redis.Client, error) {
	dsn := c.GetRedis().GetDsn()

	if dsn == "" {
		return nil, nil
	}

	cfg, err := redis.ParseURL(c.GetRedis().GetDsn())
	if err != nil {
		return nil, err
	}

	if c.GetRedis().GetMaxActiveConnections() != 0 {
		cfg.MaxActiveConns = int(c.GetRedis().GetMaxActiveConnections())
	}

	client := redis.NewClient(cfg)

	if err := redisotel.InstrumentTracing(client); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}
