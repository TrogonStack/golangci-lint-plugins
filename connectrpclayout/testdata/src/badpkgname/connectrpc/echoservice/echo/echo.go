package echoer // want "package badpkgname/connectrpc/echoservice/echo serves rpc Echo, so it must be named echo, got echoer"

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

var _ echov1connect.EchoServiceEchoHandlerFunc = Handler

func Handler(context.Context, *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{}, nil
}
