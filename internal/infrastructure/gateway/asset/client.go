package asset

import (
	"context"
	"log/slog"

	"github.com/go-kratos/kratos/v3/transport/grpc"
	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/identity/internal/conf"
)

func NewClient(
	c *conf.Gateway,
	log *slog.Logger,
) (v1.AssetServiceClient, func(), error) {
	conn, err := grpc.NewClient(
		context.Background(),
		grpc.WithEndpoint(c.Asset.Endpoint),
		grpc.WithTimeout(c.Asset.Timeout.AsDuration()),
	)
	if err != nil {
		log.Error(
			"failed to create asset service gRPC client",
			slog.String("endpoint", c.Asset.Endpoint),
			slog.Any("error", err),
		)
		return nil, nil, err
	}
	log.Info(
		"asset service gRPC client initialized",
		slog.String("endpoint", c.Asset.Endpoint),
	)

	cleanup := func() {
		_ = conn.Close()
		log.Info(
			"asset service gRPC client closed",
			slog.String("endpoint", c.Asset.Endpoint),
		)
	}

	return v1.NewAssetServiceClient(conn), cleanup, nil
}
