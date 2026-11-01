package factorycompat

import (
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Build(meter metric.Meter) {
	create := meter.Int64Counter
	_, _ = create(semconv.ServerActiveRequestsName)
	expression := metric.Meter.Int64Counter
	_, _ = expression(meter, semconv.ServerActiveRequestsName)
}
