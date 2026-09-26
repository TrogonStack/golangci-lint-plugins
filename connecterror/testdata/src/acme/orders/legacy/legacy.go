package legacy

import "github.com/bufbuild/connect-go"

func Handler() error {
	return connect.NewError(0, nil) // want "use acme/orders/rpcerr.NewError instead of connect.NewError"
}
