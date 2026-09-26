package echostream

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

func Handler(context.Context, *echov1.EchoStreamRequest) (*echov1.EchoStreamResponse, error) {
	return &echov1.EchoStreamResponse{}, nil
}

func NewHandler() echov1connect.EchoServiceEchoStreamHandlerFunc { // want "package bothforms/connectrpc/echoservice/echostream must declare Handler or NewHandler, not both"
	return Handler
}
