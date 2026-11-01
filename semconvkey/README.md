# semconvkey

Checks that OpenTelemetry attribute keys and metric names come from
[upstream semantic conventions](https://pkg.go.dev/go.opentelemetry.io/otel/semconv)
or configured semantic-convention packages.

Outside trusted packages, it reports:

- Raw-key attribute constructors such as `attribute.String("app.tier", value)`,
  including references saved or passed as function values.
- `attribute.Key` conversions and methods whose key does not come directly
  from a semantic-convention package.
- `attribute.KeyValue` literals whose `Key` has no approved origin.
- Direct metric constructors whose name has no approved origin.
- Composite literals of configured opaque types with only unexported fields.

An approved origin is a constant, variable, field, or function result declared
in `go.opentelemetry.io/otel/semconv` or a configured package, including their
subpackages. A local copy does not count. Other SDK symbols, such as
`otel.Version()`, do not define semantic keys or names.

The analyzer resolves declarations through Go types. Import aliases, dot
imports, method expressions, and embedded methods retain their identity.
Tests are skipped. Generated files are skipped by default.

## Generated instruments

Set `generated-instruments: true` when configured packages build complete
instruments rather than only declaring their names. It rejects native metric
constructors regardless of the name, including captures such as:

```go
create := meter.Int64Counter
```

It also recognizes locally declared factory methods with the exact name and
signature of an OpenTelemetry `metric.Meter` constructor. Unrelated methods
with different result types remain valid. In the default name-only mode,
saving a metric factory remains allowed.

Use the configured package's constructor instead. The rule does not prohibit
SDK/provider setup or raw instrument `Add`, `Record`, and `Observe` calls.
It does not extend to span, event, or log naming.

## Settings

| Setting                 | Meaning                                                                                                                                | Default |
| ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `allowed-packages`      | Import paths authorized to define semantic keys and names, including their subpackages. Each final path segment must end in `semconv`. | none    |
| `generated-instruments` | Require native metric construction to remain in trusted packages, including constructor references and matching local factory methods. | `false` |
| `check-generated-files` | Check generated files outside trusted packages.                                                                                        | `false` |

Upstream implementation packages under `go.opentelemetry.io` remain exempt
from analysis. This exemption does not make every upstream symbol an approved
semantic definition.

For generated-directory enforcement, enable the analyzer setting and align
**golangci-lint's issue filtering**:

```yaml
linters:
  enable:
    - semconvkey
  exclusions:
    generated: disable
  settings:
    custom:
      semconvkey:
        type: module
        settings:
          allowed-packages:
            - example.com/service/internal/appsemconv
          generated-instruments: true
          check-generated-files: true
```

Golangci-lint excludes issues from generated files by default. Without
`linters.exclusions.generated: disable`, it can hide diagnostics even when
this analyzer checks the file. Other linters may need their own path-specific
exclusions. Do not exclude this linter for untrusted generated directories.

A trusted package path or generated header does not authenticate Weaver
output. Keep a pinned generation command and a drift check that compares its
output with the committed files. The analyzer also does not prove arbitrary
runtime provenance or that a typed value is nonzero. The existing restriction
on opaque composite literals remains. Implicit zero declarations and later
constructor assignment are unchanged.

## Verification

Run the Go tests with the version declared in `go.mod`:

```sh
go test -race ./...
go vet ./...
```

Build a custom golangci-lint binary with this checkout to test issue visibility.

For a local plugin build, use `path: .` instead of the plugin's released
`version` in `.custom-gcl.yml`. Run the builder from this repository's root.

```sh
SEMCONVKEY_GOLANGCI_BINARY=/path/to/custom-binary go test -tags=integration ./semconvkey/plugin
```

The integration test is explicit and fails if the binary is not supplied.
It checks the host's generated-file filtering as well as the analyzer.
