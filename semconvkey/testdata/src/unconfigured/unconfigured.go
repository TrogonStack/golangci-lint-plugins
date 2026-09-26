// Package unconfigured is the configured fixture's calls with no
// allowed-packages: a project's own package is not trusted until it is
// listed.
package unconfigured

import (
	"example.com/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Attrs(meter metric.Meter, v string) []attribute.KeyValue {
	_, _ = meter.Int64Counter(telemetry.RequestCountName) // want `metric name must come from`
	return []attribute.KeyValue{
		telemetry.TierKey.String(v), // want `attribute key must come from`
		telemetry.Tier(v),
		attribute.Key(telemetry.TierName).String(v), // want `attribute key must come from`
		{Key: telemetry.TierKey, Value: v},          // want `attribute key must come from`
	}
}
