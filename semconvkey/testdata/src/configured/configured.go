// Package configured is a fixture for a project's own names, taken from a
// package listed in allowed-packages.
package configured

import (
	"example.com/appsemconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Attrs(meter metric.Meter, v string) []attribute.KeyValue {
	_, _ = meter.Int64Counter(appsemconv.RequestCountName)
	return []attribute.KeyValue{
		appsemconv.TierKey.String(v),
		appsemconv.Tier(v),
		attribute.Key(appsemconv.TierName).String(v),
		{Key: appsemconv.TierKey, Value: v},
	}
}
