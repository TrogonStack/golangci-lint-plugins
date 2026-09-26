// Package goodnewhandlererr is a fixture for confirming that httphandler
// reports nothing for a NewHandler that rejects its options with an error.
package goodnewhandlererr

import (
	"errors"
	"net/http"
)

type HandlerOptions struct {
	Dependency *int
}

func NewHandler(opts HandlerOptions) (http.Handler, error) {
	if opts.Dependency == nil {
		return nil, errors.New("dependency is required")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), nil
}
