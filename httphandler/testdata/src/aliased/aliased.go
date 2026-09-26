// Package aliased is a fixture for confirming that httphandler accepts a
// Handler spelled through aliases of the net/http types.
package aliased

import "aliased/web"

func Handler(w web.ResponseWriter, r web.RequestPtr) {}
