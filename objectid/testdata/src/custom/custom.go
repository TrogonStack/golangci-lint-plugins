// Package custom is a fixture for custom name-patterns.
package custom

func Get(orderRef string, orderID string, ref int) {} // want `orderRef is a raw string; use a typed ID` `ref is a raw int; use a typed ID`
