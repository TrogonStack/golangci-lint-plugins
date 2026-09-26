# connectrpclayout

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
