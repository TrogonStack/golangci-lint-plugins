package echo_stream

import (
	"context"

	"fake/gen/echov1"
	"fake/gen/echov1/echov1connect"
)

func NewHandler(prefix string, suffix string) echov1connect.EchoServiceEchoStreamHandlerFunc {
	return func(_ context.Context, req *echov1.EchoStreamRequest) (*echov1.EchoStreamResponse, error) {
		return &echov1.EchoStreamResponse{Message: prefix + req.Message + suffix}, nil
	}
}
