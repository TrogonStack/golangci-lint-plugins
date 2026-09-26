package echostream

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

type HandlerOptions struct {
	Prefix string
}

func NewHandler(opts HandlerOptions) echov1connect.EchoServiceEchoStreamHandlerFunc {
	return func(context.Context, *echov1.EchoStreamRequest) (*echov1.EchoStreamResponse, error) {
		return &echov1.EchoStreamResponse{}, nil
	}
}
