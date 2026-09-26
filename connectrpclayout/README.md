# connectrpclayout

Checks where Connect rpc handlers live, for services generated with
[protoc-gen-connect-go-servicestruct](https://github.com/TrogonStack/protoc-gen).
A package declaring a handler for one of the generated service struct's
`<Rpc>Func` fields must:

- sit at `<service>/<rpc>`, under any root directory;
- be named for the rpc, the same as its directory;
- declare the handler in the file named for the package, as `Handler` (the
  handler itself) or `NewHandler` (a constructor returning it, taking whatever
  parameters it needs), never both.

Names are spelled in lowercase by default, the way Go's package naming advice
asks, so the `EchoStream` rpc of `EchoService` lives in
`echoservice/echostream/echostream.go`. Under `snake_case` it lives in
`echo_service/echo_stream/echo_stream.go` instead, and initialisms stay one
word: `GetHTTPStatus` is `get_http_status`.

## Settings

| Setting  | Values                                 | Default      |
| -------- | -------------------------------------- | ------------ |
| `naming` | `snake_case` (`echo_service/echo_stream`), `lowercase` (`echoservice/echostream`) | `lowercase`  |

```yaml
linters:
  settings:
    custom:
      connectrpclayout:
        type: module
        settings:
          naming: snake_case
```
