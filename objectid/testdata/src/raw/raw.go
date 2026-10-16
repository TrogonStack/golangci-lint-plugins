// Package raw is a fixture for every shape objectid reports.
package raw

import "github.com/google/uuid"

type Text = string

func Param(orderID string) {} // want `orderID is a raw string; use a typed ID`

func Result() (userID int64, err error) { return 0, nil } // want `userID is a raw int64; use a typed ID`

func Pointer(orderID *string) {} // want `orderID is a raw \*string; use a typed ID`

func Slice(orderIDs []string) {} // want `orderIDs is a raw \[\]string; use a typed ID`

func Array(slotIds [4]uint32) {} // want `slotIds is a raw \[4\]uint32; use a typed ID`

func UUID(listingID uuid.UUID) {} // want `listingID is a raw uuid.UUID; use a typed ID`

func Alias(sellerId Text) {} // want `sellerId is a raw string; use a typed ID`

func Bare(id string, ids []int) {} // want `id is a raw string; use a typed ID` `ids is a raw \[\]int; use a typed ID`

func Grouped(orderID, userID string) {} // want `orderID is a raw string; use a typed ID` `userID is a raw string; use a typed ID`

type Order struct {
	ID       string // want `ID is a raw string; use a typed ID`
	BuyerIDs []*int // want `BuyerIDs is a raw \[\]\*int; use a typed ID`
	Ref      uuid.UUID
	Inner    struct {
		ParentID uint // want `ParentID is a raw uint; use a typed ID`
	}
}

type Repo interface {
	Find(orderID string) (*Order, error) // want `orderID is a raw string; use a typed ID`
}

type repo struct{}

// Find implements a same-package interface, so it is still reported: the
// signature is this package's own to change.
func (repo) Find(orderID string) (*Order, error) { return nil, nil } // want `orderID is a raw string; use a typed ID`

func Body() {
	type local struct {
		TraceID string // want `TraceID is a raw string; use a typed ID`
	}
	_ = func() {
		_ = struct {
			SpanID string // want `SpanID is a raw string; use a typed ID`
		}{}
	}
}

// OrderID embeds uuid.UUID, which is the typed ID the rule asks for.
type OrderID struct {
	uuid.UUID
}
