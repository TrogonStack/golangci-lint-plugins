// Package store declares an interface another package implements, so its
// signature is imposed on the implementation.
package store

type Getter interface {
	Get(orderID string) error // want `orderID is a raw string; use a typed ID`
}

type Lister interface {
	List(userID string) ([]string, error) // want `userID is a raw string; use a typed ID`
	Count() int
}
