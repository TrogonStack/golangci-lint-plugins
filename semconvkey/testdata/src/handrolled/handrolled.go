// Package handrolled is a fixture for instruments created outside an allowed
// package with generated-instruments set: an allowed name does not make a
// hand-rolled instrument acceptable, only the allowed package's own
// constructor does.
package handrolled

import (
	"example.com/appsemconv"

	"go.opentelemetry.io/otel/metric"
)

func Instruments(meter metric.Meter) {
	_, _ = meter.Int64Counter(appsemconv.RequestCountName) // want `instrument must be created by a package listed in allowed-packages`
	_, _ = meter.Int64Counter("app.request.count")         // want `instrument must be created by a package listed in allowed-packages`
	_, _ = appsemconv.NewRequestCount(meter)
}
