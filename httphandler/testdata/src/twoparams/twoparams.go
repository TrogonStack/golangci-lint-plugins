package twoparams

import "net/http"

type HandlerOptions struct{}

func NewHandler(opts HandlerOptions, extra int) http.Handler { // want "NewHandler must take exactly one HandlerOptions declared in this package, got 2 parameters"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
