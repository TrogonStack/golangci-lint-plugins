# golangci-lint-plugins

Custom linters for [golangci-lint](https://golangci-lint.run) using the
[module plugin system](https://golangci-lint.run/plugins/module-plugins/).

## Linters

### staticerr

Reports any `errors.New` call outside a package-level `var`, exported or not. An error built
where it is returned is a new value on every call, so `errors.Is` can never
match it. `_test.go` files are skipped.

### connectrpclayout

Checks where Connect rpc handlers live, for services generated with
[protoc-gen-connect-go-servicestruct](https://github.com/TrogonStack/protoc-gen).
A package declaring a handler for one of the generated service struct's
`<Rpc>Func` fields must:

- sit at `<service>/<rpc>`, under any root directory;
- be named for the rpc;
- declare the handler in the file named for the package, as `Handler` (the
  handler itself) or `NewHandler` (a constructor returning it, taking whatever
  parameters it needs), never both.

Case and underscores are ignored when matching names, so both
`connectrpc/echoservice/echostream/echostream.go` and
`grpc/echo_service/echo_stream/echo_stream.go` are accepted for the
`EchoStream` rpc of `EchoService`.

## Usage

Add the plugins to `.custom-gcl.yml`:

```yaml
version: v2.13.2
plugins:
  - module: github.com/TrogonStack/golangci-lint-plugins
    import: github.com/TrogonStack/golangci-lint-plugins/staticerr/plugin
    version: v0.1.0
  - module: github.com/TrogonStack/golangci-lint-plugins
    import: github.com/TrogonStack/golangci-lint-plugins/connectrpclayout/plugin
    version: v0.1.0
```

Build the custom binary with `golangci-lint custom`, then enable the linters in
`.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - staticerr
    - connectrpclayout
  settings:
    custom:
      staticerr:
        type: module
      connectrpclayout:
        type: module
```
