package handlervar

import "net/http"

var Handler http.Handler // want "Handler must be func\\(http.ResponseWriter, \\*http.Request\\), got var of type http.Handler"
