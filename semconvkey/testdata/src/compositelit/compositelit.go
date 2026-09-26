// Package compositelit is a fixture for an attribute.KeyValue literal, which
// is judged by where its Key comes from.
package compositelit

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Attrs(v string) []attribute.KeyValue {
	return []attribute.KeyValue{
		{Key: semconv.MessagingSystemKey, Value: v},
		attribute.KeyValue{Value: v, Key: semconv.MessagingSystemKey},
		{Key: "app.user.tier", Value: v}, // want `attribute key must come from`
		attribute.KeyValue{semconv.MessagingSystemKey, v},
		attribute.KeyValue{"app.user.tier", v}, // want `attribute key must come from`
		{Value: v},                             // want `attribute key must come from`
	}
}
