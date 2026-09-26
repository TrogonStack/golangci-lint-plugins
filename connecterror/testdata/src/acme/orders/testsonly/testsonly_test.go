package testsonly

import "connectrpc.com/connect"

var _ = connect.NewError(connect.CodeInvalidArgument, nil)
