package echo // want "package wrongpath/connectrpc/pingservice/echo serves rpc Echo of EchoService, so it must be at echoservice/echo"

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

var _ echov1connect.EchoServiceEchoHandlerFunc = Handler

func Handler(context.Context, *echov1.EchoRequest) (*echov1.EchoResponse, error) {
	return &echov1.EchoResponse{}, nil
}
