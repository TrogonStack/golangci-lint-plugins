// Package implements is a fixture for methods whose signature an interface
// from another package imposes.
package implements

import "example.com/store"

var (
	_ store.Getter = getter{}
	_ store.Lister = (*lister)(nil)
)

type getter struct{}

func (getter) Get(orderID string) error { return nil }

type lister struct{}

func (*lister) List(userID string) ([]string, error) { return nil, nil }

func (*lister) Count() int { return 0 }

// Delete is not part of any interface, so it is reported.
func (*lister) Delete(userID string) {} // want `userID is a raw string; use a typed ID`

// half has Lister's List but not its Count, so it does not implement
// Lister and the signature is still its own.
type half struct{}

func (half) List(userID string) ([]string, error) { return nil, nil } // want `userID is a raw string; use a typed ID`

// other has a Get whose signature differs from Getter's.
type other struct{}

func (other) Get(orderID int) error { return nil } // want `orderID is a raw int; use a typed ID`
