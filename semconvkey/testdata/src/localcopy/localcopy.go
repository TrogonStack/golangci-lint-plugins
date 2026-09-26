// Package localcopy is a fixture for a local copy of a semconv key, which is
// reported even unchanged: the rule is about where a name is named, and a
// local is a place it could be respelled.
package localcopy

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

var system = semconv.MessagingSystemKey

func System(v string) attribute.KeyValue {
	return system.String(v) // want `attribute key must come from`
}
