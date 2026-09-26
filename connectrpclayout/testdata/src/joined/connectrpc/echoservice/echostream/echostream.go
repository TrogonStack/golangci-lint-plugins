package echostream // want "package joined/connectrpc/echoservice/echostream serves rpc EchoStream of EchoService, so it must be at echo_service/echo_stream" "package joined/connectrpc/echoservice/echostream serves rpc EchoStream, so it must be named echo_stream, got echostream"

import (
	"fake/gen/echov1/echov1connect"
)

func NewHandler() echov1connect.EchoServiceEchoStreamHandlerFunc {
	return nil
}
