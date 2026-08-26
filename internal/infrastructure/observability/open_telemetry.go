package observability

import (
	"context"
	"errors"
	"fmt"

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

	if c != nil && c.Tracing != nil {
		exporter, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpointURL(c.Tracing.Endpoint),
		)
		if err != nil {
			return nil, fmt.Errorf("create OTLP trace exporter: %w", err)
		}

		opts := []sdktrace.TracerProviderOption{
			sdktrace.WithResource(res),
			sdktrace.WithBatcher(exporter),
		}

		if c.Tracing.SampleRatio != nil {
			opts = append(opts,
				sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(*c.Tracing.SampleRatio))),
			)
		}

		tracerProvider = sdktrace.NewTracerProvider(opts...)
	}

	if c != nil && c.Metrics != nil {
		exporter, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpointURL(c.Metrics.Endpoint),
		)
		if err != nil {
			if tracerProvider != nil {
				_ = tracerProvider.Shutdown(ctx)
			}

			return nil, fmt.Errorf("create OTLP metric exporter: %w", err)
		}

		reader := sdkmetric.NewPeriodicReader(
			exporter,
			sdkmetric.WithInterval(c.Metrics.ExportInterval.AsDuration()),
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
				shutdownErrors = append(shutdownErrors,
					fmt.Errorf("shutdown meter provider: %w", err),
				)
			}
		}
		if tracerProvider != nil {
			if err := tracerProvider.Shutdown(ctx); err != nil {
				shutdownErrors = append(shutdownErrors,
					fmt.Errorf("shutdown tracer provider: %w", err),
				)
			}
		}

		return errors.Join(shutdownErrors...)
	}, nil
}
