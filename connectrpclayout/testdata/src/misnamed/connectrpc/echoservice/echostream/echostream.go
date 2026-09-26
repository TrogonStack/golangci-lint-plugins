package echostream

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

func Build() echov1connect.EchoServiceEchoStreamHandlerFunc { // want "Build builds the handler for rpc EchoStream, so it must be named NewHandler"
	return Handle
}

func Handle(context.Context, *echov1.EchoStreamRequest) (*echov1.EchoStreamResponse, error) { // want "Handle is the handler for rpc EchoStream, so it must be named Handler"
	return &echov1.EchoStreamResponse{}, nil
}
