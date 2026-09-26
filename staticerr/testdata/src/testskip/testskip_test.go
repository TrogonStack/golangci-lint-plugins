package testskip

import (
	"errors"
	"testing"
)

// A one-off error standing in for a dependency's failure is a fixture, not a
// condition of this package, and there is no caller to match it. Nothing is
// reported here.
func TestValidate(t *testing.T) {
	stub := func() error { return errors.New("the database is gone") }

	if err := stub(); err == nil {
		t.Fatal("want an error")
	}

	if err := Validate(""); !errors.Is(err, ErrNoName) {
		t.Fatalf("got %v, want ErrNoName", err)
	}
}
