// Package dotimport is a fixture for the attribute constructors called
// through a dot import, where the callee is a bare identifier rather than a
// selector.
package dotimport

import (
	. "go.opentelemetry.io/otel/attribute"
)

func Attrs(v string) []KeyValue {
	return []KeyValue{
		String("app.user.tier", v),   // want `attribute.String takes the key as a raw string`
		Key("app.user.id").String(v), // want `attribute key must come from`
	}
}
