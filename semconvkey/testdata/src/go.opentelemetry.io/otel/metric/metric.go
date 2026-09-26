// Package metric stands in for go.opentelemetry.io/otel/metric, with only the
// shapes semconvkey looks at.
package metric

type Int64Counter interface{ Add(int64) }

type Float64Histogram interface{ Record(float64) }

type Meter interface {
	Int64Counter(name string, options ...int) (Int64Counter, error)
	Float64Histogram(name string, options ...int) (Float64Histogram, error)
}
