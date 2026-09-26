// Package rpcerr is the configured replacement, the one package allowed to
// call connect.NewError.
package rpcerr

import "connectrpc.com/connect"

func NewError(err error) error {
	return connect.NewError(connect.CodeInvalidArgument, err)
}
