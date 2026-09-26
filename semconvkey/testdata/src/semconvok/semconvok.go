// Package semconvok is a fixture for keys and names taken straight from
// semconv, which is always allowed and needs no configuration.
package semconvok

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Attrs(v string) []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.MessagingSystemKey.String(v),
		(semconv.MessagingSystemKey).Bool(true),
		semconv.ServerAddress(v),
	}
}

func Defined() bool {
	return semconv.MessagingSystemKey.Defined()
}

func Instruments(meter metric.Meter) {
	_, _ = meter.Int64Counter(semconv.ServerActiveRequestsName)
	_, _ = meter.Float64Histogram(semconv.ServerRequestDuration{}.Name())
}
