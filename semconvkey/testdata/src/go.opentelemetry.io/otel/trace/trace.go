// Package trace stands in for go.opentelemetry.io/otel/trace, with only the
// shapes semconvkey looks at.
package trace

import "context"

type Span interface{ End() }

type SpanStartOption interface{}

type Tracer interface {
	Start(ctx context.Context, spanName string, options ...SpanStartOption) (context.Context, Span)
}
