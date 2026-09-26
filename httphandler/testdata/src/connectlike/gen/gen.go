package gen

import "net/http"

type Request struct {
	Header http.Header
}

type Response struct{}

type CommitHandlerFunc func(*Request) (*Response, error)
