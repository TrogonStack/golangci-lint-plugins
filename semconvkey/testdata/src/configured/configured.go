// Package configured is a fixture for a project's own names, taken from a
// package listed in allowed-packages.
package configured

import (
	"example.com/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Attrs(meter metric.Meter, v string) []attribute.KeyValue {
	_, _ = meter.Int64Counter(telemetry.RequestCountName)
	return []attribute.KeyValue{
		telemetry.TierKey.String(v),
		telemetry.Tier(v),
		attribute.Key(telemetry.TierName).String(v),
		{Key: telemetry.TierKey, Value: v},
	}
}
