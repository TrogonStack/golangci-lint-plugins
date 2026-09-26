// Package aliasednew is a fixture for confirming that httphandler accepts a
// NewHandler returning an alias of http.Handler.
package aliasednew

import (
	"net/http"

	"aliased/web"
)

type HandlerOptions struct{}

func NewHandler(opts HandlerOptions) web.HTTPHandler {
	return http.HandlerFunc(func(w web.ResponseWriter, r *web.Request) {})
}
