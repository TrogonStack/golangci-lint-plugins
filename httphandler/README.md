# httphandler

Checks the shape of a plain `net/http` handler package. A package declaring
`Handler` or `NewHandler` must:

- declare exactly one of them, never both;
- declare it in the file named for the package, as `<package>.go`;
- spell `Handler` as a func `func(http.ResponseWriter, *http.Request)`, never a
  var or a type, for a handler that needs nothing from the process around it;
- spell `NewHandler` as `func(HandlerOptions) http.Handler`, or
  `func(HandlerOptions) (http.Handler, error)` when it rejects invalid options,
  with `HandlerOptions` declared in the same package, for a handler that does.

```go
package listorders

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {}
```

```go
package placeorder

import "net/http"

type HandlerOptions struct {
	Orders OrderRepository
}

func NewHandler(opts HandlerOptions) (http.Handler, error) {
	if opts.Orders == nil {
		return nil, ErrOrderRepositoryRequired
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}), nil
}
```

Only a declaration whose signature mentions a `net/http` type is checked, so a
`NewHandler` that builds a `log/slog` handler, or a Connect rpc handler spelled
in generated types (see [connectrpclayout](../connectrpclayout)), is left
alone. `package main` and `_test.go` files are skipped.

The linter has no settings.
