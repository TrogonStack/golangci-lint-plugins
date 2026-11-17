// Package spannameoff is a fixture for a span name spelled at the call site
// while span-names is off, which reports nothing.
package spannameoff

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func Spans(ctx context.Context, tracer trace.Tracer) {
	_, _ = tracer.Start(ctx, "checkout")
}
