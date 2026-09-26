// Package shared sits beside the rpc packages of a service and declares no
// handler, so none of the rules apply to it, including NewHandler of its own.
package shared

type Formatter struct{}

func NewHandler() Formatter { return Formatter{} }
