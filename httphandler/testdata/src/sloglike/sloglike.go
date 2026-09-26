// Package sloglike is a fixture for confirming that httphandler skips a
// NewHandler that builds something other than an http handler.
package sloglike

import "io"

type TextHandler struct{}

func NewHandler(w io.Writer, opts *TextHandler) *TextHandler { return opts }
