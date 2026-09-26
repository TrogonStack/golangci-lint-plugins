// Package localconst is a fixture for a key declared as a const in the
// package using it, which is a second place the name is spelled.
package localconst

import "go.opentelemetry.io/otel/attribute"

const tierKey = attribute.Key("app.user.tier") // want `attribute key must come from`

func Tier(v string) attribute.KeyValue {
	return tierKey.String(v) // want `attribute key must come from`
}
