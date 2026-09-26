// Package zeroliteral is a fixture for composite literals of types an allowed
// package declares: one with only unexported fields is reported in each shape
// a literal of it takes, and one with an exported field or no fields at all
// is not.
package zeroliteral

import (
	"example.com/appsemconv"

	"go.opentelemetry.io/otel/attribute"
)

func Attrs() []attribute.KeyValue {
	bare := appsemconv.TierAttr{}     // want `TierAttr has no exported fields`
	pointer := &appsemconv.TierAttr{} // want `TierAttr has no exported fields`
	elided := []appsemconv.TierAttr{
		{}, // want `TierAttr has no exported fields`
		appsemconv.NewTierAttr("gold"),
	}
	elidedPointer := []*appsemconv.TierAttr{
		{}, // want `TierAttr has no exported fields`
	}
	_ = appsemconv.Options{Prefix: "app"}
	_ = appsemconv.Marker{}

	return []attribute.KeyValue{bare.KeyValue(), pointer.KeyValue(), elided[0].KeyValue(), elidedPointer[0].KeyValue()}
}
