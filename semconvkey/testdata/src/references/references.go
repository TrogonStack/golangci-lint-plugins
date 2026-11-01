package references

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Capture(meter metric.Meter) {
	create := (meter.Int64Counter) // want `instrument must be created by a package listed in allowed-packages`
	_, _ = create("app.raw")
	expression := metric.Meter.Int64Counter // want `instrument must be created by a package listed in allowed-packages`
	_, _ = expression(meter, "app.raw")
	raw := attribute.String // want `takes the key as a raw string`
	_ = raw("app.raw", "value")
	consume((meter.Float64Histogram)) // want `instrument must be created by a package listed in allowed-packages`
	fields := struct {
		create func(string, ...int) (metric.Int64Counter, error)
	}{
		create: meter.Int64Counter, // want `instrument must be created by a package listed in allowed-packages`
	}
	_, _ = fields.create("app.raw")
	bound := semconv.MessagingSystemKey.String
	_ = bound("value")
	_ = semconv.ServerAddress
	_ = attribute.Key.String // want `attribute key must come from`
}

func Return(meter metric.Meter) func(string, ...int) (metric.Int64Counter, error) {
	return meter.Int64Counter // want `instrument must be created by a package listed in allowed-packages`
}

func consume(func(string, ...int) (metric.Float64Histogram, error)) {}
