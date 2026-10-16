// Package untypedoff is untyped's conversions with untyped-constants off.
package untypedoff

type OrderID string

func Use(id OrderID) {}

func Calls() {
	_ = OrderID("abc")
	Use("abc")
}
