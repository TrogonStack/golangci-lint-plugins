// Package metricname is a fixture for a metric instrument whose name is
// spelled at the call site. The Meter call itself is fine; the name is what
// is judged.
package metricname

import "go.opentelemetry.io/otel/metric"

const requestCountName = "app.request.count"

type recorder struct{}

func (recorder) Int64Counter(name string) {}

func Instruments(meter metric.Meter, r recorder) {
	_, _ = meter.Int64Counter("app.request.count")  // want `metric name must come from`
	_, _ = meter.Float64Histogram(requestCountName) // want `metric name must come from`
	r.Int64Counter("not.an.instrument")
}
