package upstreamorigin

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Attrs() {
	_ = attribute.Key(otel.Version()).String("value")          // want `attribute key must come from`
	_ = attribute.KeyValue{Key: attribute.Key(otel.Version())} // want `attribute key must come from` `attribute key must come from`
	_ = semconv.MessagingSystemKey.String("value")
}

func MeterName(meter metric.Meter) {
	_, _ = meter.Int64Counter(otel.Version()) // want `metric name must come from`
}
