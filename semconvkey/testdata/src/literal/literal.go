// Package literal is a fixture for the attribute constructors that take the
// key as a raw string, which are reported whatever the string is.
package literal

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Attrs(v string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("app.user.tier", v),                   // want `attribute.String takes the key as a raw string`
		attribute.Bool("app.user.active", true),                // want `attribute.Bool takes the key as a raw string`
		attribute.Int64(string(semconv.MessagingSystemKey), 1), // want `attribute.Int64 takes the key as a raw string`
	}
}
