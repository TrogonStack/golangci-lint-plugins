package badoptsname

import "net/http"

type Config struct{}

func NewHandler(cfg Config) http.Handler { // want "the argument to NewHandler must be a HandlerOptions declared in this package, got Config"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
