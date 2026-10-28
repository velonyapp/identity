package asset

import (
	"context"

	"github.com/go-kratos/kratos/v3/transport/grpc"
	v1 "github.com/velonyapp/asset/gen/api/v1"
	"github.com/velonyapp/identity/internal/conf"
)

func NewClient(c *conf.Gateway) (v1.AssetServiceClient, func(), error) {
	conn, err := grpc.NewClient(
		context.Background(),
		grpc.WithEndpoint(c.Asset.Endpoint),
		grpc.WithTimeout(c.Asset.Timeout.AsDuration()),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = conn.Close()
	}

	return v1.NewAssetServiceClient(conn), cleanup, nil
}
