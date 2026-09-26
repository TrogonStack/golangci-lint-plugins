// Package autherr is a second configured replacement.
package autherr

import "connectrpc.com/connect"

func NewError(err error) error {
	return connect.NewError(connect.CodeInvalidArgument, err)
}
