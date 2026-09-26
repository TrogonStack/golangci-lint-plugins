// Package echov1connect stands in for the code
// protoc-gen-connect-go-servicestruct generates, for connectrpclayout's own
// tests only.
package echov1connect

import (
	"context"

	"fake/gen/echov1"
)

type EchoServiceEchoHandlerFunc func(context.Context, *echov1.EchoRequest) (*echov1.EchoResponse, error)

type EchoServiceEchoStreamHandlerFunc func(context.Context, *echov1.EchoStreamRequest) (*echov1.EchoStreamResponse, error)

type EchoServiceStruct struct {
	EchoFunc       EchoServiceEchoHandlerFunc
	EchoStreamFunc EchoServiceEchoStreamHandlerFunc
}
