package several

import "connectrpc.com/connect"

func Handler() error {
	return connect.NewError(connect.CodeInvalidArgument, nil) // want "use one of acme/orders/rpcerr.NewError, acme/orders/autherr.NewError instead of connect.NewError"
}
