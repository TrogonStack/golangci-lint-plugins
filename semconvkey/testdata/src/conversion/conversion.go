// Package conversion is a fixture for a conversion to attribute.Key, which is
// fine when its operand already comes from semconv and reported otherwise.
package conversion

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

func Keys(name string) []attribute.Key {
	return []attribute.Key{
		attribute.Key("app.user.tier"), // want `attribute key must come from`
		attribute.Key(name),            // want `attribute key must come from`
		attribute.Key(semconv.MessagingSystemKey),
		attribute.Key(semconv.ServerActiveRequestsName),
	}
}
