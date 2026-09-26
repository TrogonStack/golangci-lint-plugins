package badsig

import "net/http"

func Handler() http.Handler { return nil } // want "Handler must be func\\(http.ResponseWriter, \\*http.Request\\), got func\\(\\) http.Handler"
