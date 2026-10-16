// Package named is a fixture for what objectid leaves alone.
package named

import "github.com/google/uuid"

type OrderID string

type UserIDs []OrderID

type Wrapper struct{ value uuid.UUID }

func Named(orderID OrderID, userIDs UserIDs, listingID *Wrapper) {}

func NotAnID(valid bool, Grid string, paid int, name string) {}

func NoName(string, int) (string, error) { return "", nil }

func Blank(_ string) {}

func Mixed(orderID float64, userID []byte) {} // want `userID is a raw \[\]byte; use a typed ID`

func Generic[T ~string](orderID T) {}

func (orderID OrderID) Receiver() {}

func Callbacks() {
	register(func(orderID string) error { return nil })
	handler := func(userID int) {}
	handler(1)
}

func register(func(string) error) {}

type Embedded struct {
	OrderID
	Name string
}

// Literals are left alone while untyped-constants is off, the default.
func Literals() {
	_ = OrderID("abc")
	Named("abc", nil, nil)
}
