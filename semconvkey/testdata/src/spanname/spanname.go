// Package spanname is a fixture for a span whose name is spelled at the call
// site, with span-names enabled.
package spanname

import (
	"context"

	"example.com/appsemconv"
	"go.opentelemetry.io/otel/trace"
)

const checkoutSpanName = "checkout"

type tracer struct{}

func (tracer) Start(ctx context.Context, spanName string, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	return ctx, nil
}

type recorder struct{}

func (recorder) Start(ctx context.Context, name string) {}

func Spans(ctx context.Context, upstream trace.Tracer, own tracer, r recorder, route string) {
	_, _ = upstream.Start(ctx, "checkout")       // want `span name must come from`
	_, _ = upstream.Start(ctx, checkoutSpanName) // want `span name must come from`
	_, _ = upstream.Start(ctx, "send "+route)    // want `span name must come from`
	_, _ = own.Start(ctx, "checkout")            // want `span name must come from`
	_, _ = upstream.Start(ctx, appsemconv.CheckoutSpanName)
	_, _ = upstream.Start(ctx, appsemconv.SendSpanName(route))
	r.Start(ctx, "not.a.span")
}
