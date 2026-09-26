package badreturn

import "net/http"

type HandlerOptions struct{}

func NewHandler(opts HandlerOptions) (http.Handler, error) { // want "NewHandler must return http.Handler, got \\(http.Handler, error\\)"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), nil
}
