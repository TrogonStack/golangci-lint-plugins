// Package localiface is a fixture for confirming that httphandler resolves
// Handler through go/types rather than by the identifier's spelling: this
// package's Handler is a locally declared interface, not net/http's own, and
// is reported as a bad Handler for that reason.
package localiface

import "net/http"

type Handler interface { // want "Handler must be func\\(http.ResponseWriter, \\*http.Request\\), got a type declaration, not a func"
	ServeHTTP(http.ResponseWriter, *http.Request)
}
