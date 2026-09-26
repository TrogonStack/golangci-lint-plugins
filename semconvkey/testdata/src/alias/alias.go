// Package alias is a fixture for aliases of attribute.Key and
// attribute.KeyValue, which are held to the same rules as the types they
// stand for.
package alias

import "go.opentelemetry.io/otel/attribute"

type (
	K  = attribute.Key
	KV = attribute.KeyValue
)

func Attrs(v string) []attribute.KeyValue {
	return []attribute.KeyValue{
		K("app.user.tier").String(v),     // want `attribute key must come from`
		KV{Key: "app.user.id", Value: v}, // want `attribute key must come from`
	}
}
