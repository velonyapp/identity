package observability

import (
	"fmt"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const instrumentationName = "github.com/velony-app/identity/internal/infrastructure/observability"

type Metrics struct {
	Requests metric.Int64Counter
	Seconds  metric.Float64Histogram
}

func NewMetrics() (*Metrics, error) {
	meter := otel.Meter(instrumentationName)

	Requests, err := metrics.DefaultRequestsCounter(
		meter,
		metrics.DefaultServerRequestsCounterName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create server Requests counter: %w",
			err,
		)
	}

	Seconds, err := metrics.DefaultSecondsHistogram(
		meter,
		metrics.DefaultServerSecondsHistogramName,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create server duration histogram: %w",
			err,
		)
	}

	return &Metrics{
		Requests: Requests,
		Seconds:  Seconds,
	}, nil
}
