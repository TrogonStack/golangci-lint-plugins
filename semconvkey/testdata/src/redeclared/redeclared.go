package redeclared

import "go.opentelemetry.io/otel/metric"

type Counter = metric.Int64Counter

type Meter interface {
	Int64Counter(string, ...int) (Counter, error)
}

func Build(meter Meter) {
	_, _ = meter.Int64Counter("app.raw") // want `instrument must be created by a package listed in allowed-packages`
	create := meter.Int64Counter         // want `instrument must be created by a package listed in allowed-packages`
	_, _ = create("app.raw")
	expression := Meter.Int64Counter // want `instrument must be created by a package listed in allowed-packages`
	_, _ = expression(meter, "app.raw")
}

type unrelated interface {
	Int64Counter(string, ...int) (int, error)
}

func Unrelated(meter unrelated) {
	_, _ = meter.Int64Counter("app.unrelated")
	_ = meter.Int64Counter
}

func Observe(counter metric.Int64Counter) { counter.Add(1) }
