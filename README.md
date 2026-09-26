# golangci-lint-plugins

Custom linters for [golangci-lint](https://golangci-lint.run) using the
[module plugin system](https://golangci-lint.run/plugins/module-plugins/).

## Linters

### staticerr

Reports any `errors.New` call outside a package-level `var`, exported or not. An error built
where it is returned is a new value on every call, so `errors.Is` can never
match it. `_test.go` files are skipped.

## Usage

Add the plugin to `.custom-gcl.yml`:

```yaml
version: v2.13.2
plugins:
  - module: github.com/TrogonStack/golangci-lint-plugins
    import: github.com/TrogonStack/golangci-lint-plugins/staticerr/plugin
    version: v0.1.0
```

Build the custom binary with `golangci-lint custom`, then enable the linter in
`.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - staticerr
  settings:
    custom:
      staticerr:
        type: module
```
