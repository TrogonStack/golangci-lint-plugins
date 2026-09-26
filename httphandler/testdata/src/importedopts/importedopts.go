package importedopts

import (
	"net/http"

	"importedopts/otheropts"
)

func NewHandler(opts otheropts.HandlerOptions) http.Handler { // want "the argument to NewHandler must be a HandlerOptions declared in this package, got otheropts.HandlerOptions"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
