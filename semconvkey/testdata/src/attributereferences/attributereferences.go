package attributereferences

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Capture() {
	create := attribute.String // want `takes the key as a raw string`
	_ = create("app.raw", "value")
	consume(attribute.Bool)           // want `takes the key as a raw string`
	expression := attribute.Key.Int64 // want `attribute key must come from`
	_ = expression("app.raw", 1)
	copied := semconv.MessagingSystemKey
	_ = copied.String // want `attribute key must come from`
	bound := (semconv.MessagingSystemKey).String
	_ = bound("value")
	_ = (attribute.Key(semconv.MessagingSystemKey)).String
}

func consume(func(string, bool) attribute.KeyValue) {}
