# connecterror

Reports any use of Connect's `connect.NewError` outside the function a
codebase has chosen to build its Connect errors. Connect sends an error's
`Error()` text to the client as is, so a direct `connect.NewError` skips
whatever that function decides a client may read about a failure.

Every reference to `connect.NewError` is reported, whether it's called, reached
through a renamed or dot import, or passed around as a value. Both
`connectrpc.com/connect` and its earlier path `github.com/bufbuild/connect-go`
are covered. The package declaring the replacement may call
`connect.NewError`, since it has to. Generated files and `_test.go` files are
skipped.

```go
return connect.NewError(connect.CodeNotFound, ErrOrderNotFound) // reported
return rpcerr.NewError(ErrOrderNotFound)
```

## Settings

| Setting       | Value                                                        | Default  |
| ------------- | ------------------------------------------------------------ | -------- |
| `replacement` | the function to use instead, as `<import path>.<Name>`       | required |

```yaml
linters:
  settings:
    custom:
      connecterror:
        type: module
        settings:
          replacement: github.com/acme/orders/internal/rpcerr.NewError
```
