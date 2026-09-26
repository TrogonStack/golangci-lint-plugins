// Package server hands an rpc's handler on under a name of its own. It
// declares no handler, so it is not an rpc package and none of the rules
// apply to it.
package server

import (
	"fake/gen/echov1/echov1connect"

	"wiring/connectrpc/echoservice/echostream"
)

func EchoStream() echov1connect.EchoServiceEchoStreamHandlerFunc {
	return echostream.NewHandler(echostream.HandlerOptions{Prefix: "> "})
}
