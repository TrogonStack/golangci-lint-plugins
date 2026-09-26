// Package good is a fixture for confirming that httphandler reports nothing
// for a compliant Handler package.
package good

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {}
