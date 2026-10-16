// Package untyped is a fixture for untyped-constants.
package untyped

type OrderID string

type WorkID int

type Kind string

const literal = "abc"

const Typed OrderID = "well-known"

const (
	First WorkID = iota + 1
	Second
)

func Use(id OrderID) {}

func Calls(id OrderID, s string) {
	_ = OrderID("abc") // want `untyped constant converted to OrderID; build the ID with its constructor or parser instead of a literal`
	Use("abc")         // want `untyped constant converted to OrderID`
	Use(literal)       // want `untyped constant converted to OrderID`
	Use(("a" + "b"))   // want `untyped constant converted to OrderID`
	_ = WorkID(1)      // want `untyped constant converted to WorkID`
	_ = []OrderID{"x"} // want `untyped constant converted to OrderID`
	_ = id == ""
	_ = WorkID(0)
	_ = Kind("k")
	Use(Typed)
	Use(OrderID(s))
	var w WorkID = Second
	_ = w
}
