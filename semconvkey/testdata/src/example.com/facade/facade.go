package facade

import "go.opentelemetry.io/otel/metric"

type Meter interface {
	Int64Counter(string, ...int) (metric.Int64Counter, error)
}
