package bothforms

import "net/http"

type HandlerOptions struct{}

func Handler(w http.ResponseWriter, r *http.Request) {}

func NewHandler(opts HandlerOptions) http.Handler { return http.HandlerFunc(Handler) } // want "bothforms declares both Handler and NewHandler; a package has one answer to what its handler is, so declare Handler for a handler that needs nothing from the process or NewHandler for one that does"
