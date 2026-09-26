package badreturntype

import "net/http"

type HandlerOptions struct{}

func NewHandler(opts HandlerOptions) http.HandlerFunc { // want "NewHandler must return http.Handler or \\(http.Handler, error\\), got http.HandlerFunc"
	return func(w http.ResponseWriter, r *http.Request) {}
}
