// Package echov1 stands in for the messages protoc-gen-go generates, for
// connectrpclayout's own tests only.
package echov1

type EchoRequest struct{ Message string }

type EchoResponse struct{ Message string }

type EchoStreamRequest struct{ Message string }

type EchoStreamResponse struct{ Message string }
