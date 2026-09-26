package echostream // want "package mixed/grpc/echo_service/echo_stream serves rpc EchoStream, so it must be named echo_stream, got echostream"

import (
	"fake/gen/echov1/echov1connect"
)

func NewHandler() echov1connect.EchoServiceEchoStreamHandlerFunc {
	return nil
}
