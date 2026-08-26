package observability

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/velony-app/identity/internal/conf"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type ShutdownFunc func(context.Context) error

func NewOpenTelemetry(
	ctx context.Context,
	c *conf.Observability,
	serviceName string,
	serviceVersion string,
	serviceInstanceID string,
) (ShutdownFunc, error) {
	res, err := resource.New(
		ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("service.version", serviceVersion),
			attribute.String("service.instance.id", serviceInstanceID),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, fmt.Errorf("create OpenTelemetry resource: %w", err)
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	var tracerProvider *sdktrace.TracerProvider
	var meterProvider *sdkmetric.MeterProvider

	if c != nil && c.Tracing != nil && c.Tracing.Endpoint != "" {
		if err := validateEndpoint(c.Tracing.Endpoint); err != nil {
			return nil, fmt.Errorf("invalid tracing endpoint: %w", err)
		}

		if *c.Tracing.SampleRatio < 0 || *c.Tracing.SampleRatio > 1 {
			return nil, fmt.Errorf(
				"tracing sample ratio must be between 0 and 1: %f",
				c.Tracing.SampleRatio,
			)
		}

		exporter, err := otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpointURL(c.Tracing.Endpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
		}

		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),

			sdktrace.WithSampler(
				sdktrace.ParentBased(
					sdktrace.TraceIDRatioBased(
						*c.Tracing.SampleRatio,
					),
				),
			),

			sdktrace.WithBatcher(exporter),
		)
	}

	if c != nil && c.Metrics != nil && c.Metrics.Endpoint != "" {
		if err := validateEndpoint(c.Metrics.Endpoint); err != nil {
			if tracerProvider != nil {
				_ = tracerProvider.Shutdown(ctx)
			}

			return nil, fmt.Errorf("invalid metrics endpoint: %w", err)
		}

		exporter, err := otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithEndpointURL(c.Metrics.Endpoint),
		)
		if err != nil {
			if tracerProvider != nil {
				_ = tracerProvider.Shutdown(ctx)
			}

			return nil, fmt.Errorf("create OTLP metric exporter: %w", err)
		}

		readerOptions := make(
			[]sdkmetric.PeriodicReaderOption,
			0,
			1,
		)

		if c.Metrics.ExportInterval != nil {
			interval := c.Metrics.ExportInterval.AsDuration()

			if interval > 0 {
				readerOptions = append(
					readerOptions,
					sdkmetric.WithInterval(interval),
				)
			}
		}

		reader := sdkmetric.NewPeriodicReader(
			exporter,
			readerOptions...,
		)

		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(reader),
		)
	}

	if tracerProvider != nil {
		otel.SetTracerProvider(tracerProvider)
	}
	if meterProvider != nil {
		otel.SetMeterProvider(meterProvider)
	}

	return func(ctx context.Context) error {
		var shutdownErrors []error

		if meterProvider != nil {
			if err := meterProvider.Shutdown(ctx); err != nil {
				shutdownErrors = append(
					shutdownErrors,
					fmt.Errorf("shutdown meter provider: %w", err),
				)
			}
		}
		if tracerProvider != nil {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				shutdownErrors = append(
					shutdownErrors,
					fmt.Errorf("shutdown tracer provider: %w", err),
				)
			}
		}

		return errors.Join(shutdownErrors...)
	}, nil
}

func validateEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf(
			"endpoint must use http or https scheme",
		)
	}

	if u.Host == "" {
		return fmt.Errorf("endpoint must contain a host")
	}

	return nil
}
