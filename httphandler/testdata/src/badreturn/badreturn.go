package badreturn

import "net/http"

type HandlerOptions struct{}

func NewHandler(opts HandlerOptions) (error, http.Handler) { // want "NewHandler must return http.Handler or \\(http.Handler, error\\), got \\(error, http.Handler\\)"
	return nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
