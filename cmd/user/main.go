package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/velony-app/identity/internal/conf"
	"github.com/velony-app/identity/internal/infrastructure/observability"

	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3"
	"github.com/go-kratos/kratos/v3/config"
	"github.com/go-kratos/kratos/v3/config/env"
	"github.com/go-kratos/kratos/v3/config/file"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/go-kratos/kratos/v3/transport/http"

	_ "go.uber.org/automaxprocs"
)

var (
	Name = "velony-identity"

	Version = "dev"

	InstanceID, _ = os.Hostname()

	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger *slog.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(InstanceID),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func newLogger(
	c *conf.Observability_Logging,
) (*slog.Logger, error) {
	level := slog.LevelInfo

	if c != nil && c.Level != "" {
		if err := level.UnmarshalText([]byte(c.Level)); err != nil {
			return nil, fmt.Errorf("parse log level %q: %w", c.Level, err)
		}
	}

	logger := log.NewLogger(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}),
		log.WithExtractor(tracing.TraceAttrs),
	).With(
		slog.String("service.InstanceID", InstanceID),
		slog.String("service.name", Name),
		slog.String("service.version", Version),
	)

	return logger, nil
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			env.NewSource("VELONY_PROFILE_"),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger, err := newLogger(bc.Observability.GetLogging())
	if err != nil {
		panic(err)
	}

	log.SetDefault(logger)

	ctx := context.Background()

	shutdownOTel, err := observability.NewOpenTelemetry(
		ctx,
		bc.Observability,
		Name,
		Version,
		InstanceID,
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := shutdownOTel(ctx); err != nil {
			logger.Error(
				"failed to shutdown OpenTelemetry",
				"error", err,
			)
		}
	}()

	app, cleanup, err := wireApp(
		bc.Data,
		bc.Transport,
		bc.Observability,
		bc.Auth,
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
