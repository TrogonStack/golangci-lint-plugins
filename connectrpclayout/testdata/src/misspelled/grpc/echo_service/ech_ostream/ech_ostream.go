package ech_ostream // want "package misspelled/grpc/echo_service/ech_ostream serves rpc EchoStream of EchoService, so it must be at echo_service/echo_stream" "package misspelled/grpc/echo_service/ech_ostream serves rpc EchoStream, so it must be named echo_stream, got ech_ostream"

import (
	"fake/gen/echov1/echov1connect"
)

func NewHandler() echov1connect.EchoServiceEchoStreamHandlerFunc {
	return nil
}
