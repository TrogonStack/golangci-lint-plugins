# connecterror

Reports any use of Connect's `connect.NewError` outside the functions a
codebase has chosen to build its Connect errors. Connect sends an error's
`Error()` text to the client as is, so a direct `connect.NewError` skips
whatever those functions decide a client may read about a failure.

Every reference to `connect.NewError` is reported, whether it's called, reached
through a renamed or dot import, or passed around as a value. Both
`connectrpc.com/connect` and its earlier path `github.com/bufbuild/connect-go`
are covered. A package declaring a replacement may call
`connect.NewError`, since it has to. Generated files and `_test.go` files are
skipped.

```go
return connect.NewError(connect.CodeNotFound, ErrOrderNotFound) // reported
return rpcerr.NewError(connect.CodeNotFound, ErrOrderNotFound)
```

A replacement's signature is its own business: it can mirror
`connect.NewError`, take only the error, or anything else.

## Settings

| Setting        | Value                                                        | Default  |
| -------------- | ------------------------------------------------------------ | -------- |
| `replacements` | the functions to use instead, each as `<import path>.<Name>` | required |

```yaml
linters:
  settings:
    custom:
      connecterror:
        type: module
        settings:
          replacements:
            - github.com/acme/orders/internal/rpcerr.NewError
            - github.com/acme/orders/internal/autherr.NewError
```
