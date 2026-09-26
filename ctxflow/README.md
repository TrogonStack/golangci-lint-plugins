# ctxflow

Checks that a `context.Context` keeps flowing down every call instead of
being dropped. It reports:

- a call to `F` when `F` takes no context but its package or receiver declares
  `FContext` or `FWithContext` taking a `context.Context` first and otherwise
  the same parameters and results, such as `slog.Info`, `logger.Warn`,
  `db.Query`, `exec.Command` or `http.NewRequest`. The sibling must be
  declared on the same type and callable on the same operand, so a method
  promoted from an embedded field, or a pointer-receiver sibling of a call on
  a value that is not addressable, is not suggested;
- any use of `context.Background` or `context.TODO` outside `func main` in
  `package main`, a `func init`, `TestMain`, or an `Example` function. A
  function literal inside one of those is still checked, since it can capture
  the context its enclosing function built.

A log record written without the context is not attached to the trace it
happened under, and a query or request sent without it outlives the caller
that gave up on it. The sibling rule follows the shape of an API instead of a
list of loggers, so any library using the `FContext` convention is covered.

```go
func (s *Service) Charge(ctx context.Context, id OrderID) error {
	slog.Info("charging")             // reported
	slog.InfoContext(ctx, "charging") // ok

	go s.audit(context.Background(), id)       // reported
	go s.audit(context.WithoutCancel(ctx), id) // ok

	_, err := s.db.Exec(chargeQuery, id)            // reported
	_, err = s.db.ExecContext(ctx, chargeQuery, id) // ok
	return err
}
```

Where no context seems available, there usually is one:

| Situation                                   | Use                                   |
| ------------------------------------------- | ------------------------------------- |
| A test                                      | `t.Context()`                         |
| Work that must outlive the caller's cancellation, such as a shutdown flush | `context.WithoutCancel(ctx)` |
| `main` delegating to `run()`                | Build the root in `main` and pass it to `run(ctx)` |
| A library callback that receives no context | `//nolint:ctxflow // <reason>`        |

Generated files are skipped. The linter has no settings.
