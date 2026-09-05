package observability

import (
	"context"
	"errors"
	"time"

	"github.com/velonyapp/identity/internal/conf"
	"github.com/velonyapp/identity/internal/info"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type OpenTelemetry struct {
	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider
}

func NewOpenTelemetry(ctx context.Context, c *conf.Observability, i *info.Service) (*OpenTelemetry, func(), error) {
	if c == nil {
		return &OpenTelemetry{}, func() {}, nil
	}

	res, err := resource.New(ctx,
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			attribute.String("service.name", i.Name),
			attribute.String("service.version", i.Version),
			attribute.String("service.instance.id", i.InstanceID),
		),
		resource.WithFromEnv(),
	)
	if err != nil {
		return nil, nil, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	o := &OpenTelemetry{}

	shutdown := func(ctx context.Context) error {
		var errs []error

		if o.meterProvider != nil {
			if err := o.meterProvider.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
		}

		if o.tracerProvider != nil {
			if err := o.tracerProvider.Shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Join(errs...)
	}

	if c.Tracing != nil {
		exporter, err := otlptracegrpc.New(
			ctx,
			otlptracegrpc.WithEndpointURL(c.Tracing.Endpoint),
		)
		if err != nil {
			return nil, nil, err
		}

		opts := []trace.TracerProviderOption{
			trace.WithResource(res),
			trace.WithBatcher(exporter),
		}

		if c.Tracing.SampleRatio != nil {
			opts = append(
				opts,
				trace.WithSampler(
					trace.ParentBased(
						trace.TraceIDRatioBased(*c.Tracing.SampleRatio),
					),
				),
			)
		}

		o.tracerProvider = trace.NewTracerProvider(opts...)
	}

	if c.Metrics != nil {
		exporter, err := otlpmetricgrpc.New(
			ctx,
			otlpmetricgrpc.WithEndpointURL(c.Metrics.Endpoint),
		)
		if err != nil {
			_ = shutdown(ctx)

			return nil, nil, err
		}

		reader := metric.NewPeriodicReader(
			exporter,
			metric.WithInterval(c.Metrics.ExportInterval.AsDuration()),
		)

		o.meterProvider = metric.NewMeterProvider(
			metric.WithResource(res),
			metric.WithReader(reader),
		)
	}

	if o.tracerProvider != nil {
		otel.SetTracerProvider(o.tracerProvider)
	}

	if o.meterProvider != nil {
		otel.SetMeterProvider(o.meterProvider)
	}

	cleanup := func() {
		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = shutdown(ctx)
	}

	return o, cleanup, nil
}
