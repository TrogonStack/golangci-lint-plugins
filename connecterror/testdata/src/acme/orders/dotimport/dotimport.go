package dotimport

import . "connectrpc.com/connect"

func Handler() error {
	return NewError(CodeInvalidArgument, nil) // want "use acme/orders/rpcerr.NewError instead of connect.NewError"
}
