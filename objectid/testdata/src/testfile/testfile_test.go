package testfile

import "testing"

func lookup(orderID string) string { return orderID }

func TestLookup(t *testing.T) {
	if lookup("o") != "o" {
		t.Fatal()
	}
}
