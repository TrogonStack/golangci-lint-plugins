package echo

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

var _ echov1connect.EchoServiceEchoHandlerFunc = Handler

func Handler(context.Context, *echov1.EchoRequest) (*echov1.EchoResponse, error) { // want "Handler must be declared in echo.go, not in other.go"
	return &echov1.EchoResponse{}, nil
}
