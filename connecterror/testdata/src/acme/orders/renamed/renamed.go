package renamed

import rpc "connectrpc.com/connect"

func Handler() error {
	return rpc.NewError(rpc.CodeInvalidArgument, nil) // want "use acme/orders/rpcerr.NewError instead of connect.NewError"
}
