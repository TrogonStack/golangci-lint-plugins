# golangci-lint-plugins

Custom linters for [golangci-lint](https://golangci-lint.run) using the
[module plugin system](https://golangci-lint.run/plugins/module-plugins/).

Each linter lives in its own top-level directory, with a `README.md` that
describes what it checks. The directory name is the linter name.

## Usage

Add a linter's plugin to `.custom-gcl.yml`, replacing `<linter>` with its
directory name:

```yaml
version: v2.13.2
plugins:
  - module: github.com/TrogonStack/golangci-lint-plugins
    import: github.com/TrogonStack/golangci-lint-plugins/<linter>/plugin
    version: v0.1.0
```

Build the custom binary with `golangci-lint custom`, then enable the linter in
`.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - <linter>
  settings:
    custom:
      <linter>:
        type: module
```
