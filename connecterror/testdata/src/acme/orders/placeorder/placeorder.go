package placeorder

import (
	"errors"

	"connectrpc.com/connect"
)

var ErrOutOfStock = errors.New("out of stock")

func Handler() error {
	return connect.NewError(connect.CodeInvalidArgument, ErrOutOfStock) // want "use acme/orders/rpcerr.NewError instead of connect.NewError"
}

func Method(e *connect.Error) {
	e.NewError()
}
