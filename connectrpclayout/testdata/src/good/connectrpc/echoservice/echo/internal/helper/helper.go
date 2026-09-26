// Package helper sits under an rpc package and declares no handler, so none
// of the rules apply to it.
package helper

func Trim(s string) string { return s }
