// Package goodnewhandler is a fixture for confirming that httphandler
// reports nothing for a compliant NewHandler package.
package goodnewhandler

import "net/http"

type HandlerOptions struct {
	Dependency int
}

func NewHandler(opts HandlerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
