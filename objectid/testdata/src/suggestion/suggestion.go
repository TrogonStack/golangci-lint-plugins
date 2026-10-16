// Package suggestion is a fixture for the suggestion setting, which every
// message carries.
package suggestion

type OrderID string

func Get(orderID string) {} // want `orderID is a raw string; use a typed ID \(internal/boid\)`

func Literal() {
	_ = OrderID("abc") // want `untyped constant converted to OrderID; build the ID with its constructor or parser instead of a literal \(internal/boid\)`
}
