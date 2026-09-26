package echo

import (
	"fake/gen/echov1/echov1connect"
)

func Stub() echov1connect.EchoServiceEchoHandlerFunc { return Handler }
