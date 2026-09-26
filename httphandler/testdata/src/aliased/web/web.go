package web

import "net/http"

type (
	ResponseWriter = http.ResponseWriter
	Request        = http.Request
	HTTPHandler    = http.Handler
	RequestPtr     = *http.Request
)
