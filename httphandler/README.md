# httphandler

Checks the shape of a plain `net/http` handler package. A package declaring
`Handler` or `NewHandler` must:

- declare exactly one of them, never both;
- declare it in the file named for the package, as `<package>.go`;
- spell `Handler` as a func `func(http.ResponseWriter, *http.Request)`, never a
  var or a type, for a handler that needs nothing from the process around it;
- spell `NewHandler` as `func(HandlerOptions) http.Handler`, with
  `HandlerOptions` declared in the same package, for a handler that does.

```go
package healthz

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {}
```

```go
package avatar

import "net/http"

type HandlerOptions struct {
	Store Store
}

func NewHandler(opts HandlerOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}
```

Only a declaration whose signature mentions a `net/http` type is checked, so a
`NewHandler` that builds a `log/slog` handler, or a Connect rpc handler spelled
in generated types (see [connectrpclayout](../connectrpclayout)), is left
alone. `package main` and `_test.go` files are skipped.

The linter has no settings.
