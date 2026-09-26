# semconvkey

Checks that an OpenTelemetry attribute key or metric instrument name comes
from [semconv](https://pkg.go.dev/go.opentelemetry.io/otel/semconv) or a
package allowed to declare one, instead of being spelled at the call site.
Outside an allowed package, it reports:

- `attribute.String("app.user.tier", v)` and the other constructors taking the
  key as a raw string;
- `key.String(v)` and the other `attribute.Key` methods, unless `key` comes
  from an allowed package;
- `attribute.Key(x)`, unless `x` comes from an allowed package;
- an `attribute.KeyValue{...}` literal, unless its `Key` comes from an allowed
  package;
- `meter.Int64Counter(name)` and the other `metric.Meter` instrument methods,
  unless `name` comes from an allowed package.

A value comes from an allowed package when it names a constant, variable,
field or function result that package declares. A local copy does not count.
`*_test.go` and generated files are skipped.

## Settings

| Setting            | Values                                                         | Default |
| ------------------ | -------------------------------------------------------------- | ------- |
| `allowed-packages` | Import paths that, with every package under them, may declare keys and names. Each must name a package ending in `semconv`, such as `appsemconv` | none    |

`go.opentelemetry.io` is always allowed. An entry whose last path segment does not end in `semconv` fails plugin load, so the package that declares a project's names reads as the project's own semconv wherever it is imported.

```yaml
linters:
  settings:
    custom:
      semconvkey:
        type: module
        settings:
          allowed-packages:
            - example.com/app/internal/appsemconv
```
