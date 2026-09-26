package funcvalue

import "connectrpc.com/connect"

var newError = connect.NewError // want "use acme/orders/rpcerr.NewError instead of connect.NewError"

func Handler() error {
	return newError(connect.CodeInvalidArgument, nil)
}
