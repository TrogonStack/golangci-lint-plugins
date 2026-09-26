// Package methodexpr is a fixture for method expressions, which pass the
// receiver as the first argument: that argument is the key or the meter, and
// is judged the same as the receiver of an ordinary method call.
package methodexpr

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Attrs(v string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Key.String(semconv.MessagingSystemKey, v),
		attribute.Key.String("app.user.tier", v), // want `attribute key must come from`
	}
}

func Instruments(meter metric.Meter) {
	_, _ = metric.Meter.Int64Counter(meter, "app.request.count") // want `metric name must come from`
}
