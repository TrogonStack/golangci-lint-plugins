// Package connectlike is a fixture for confirming that httphandler skips a
// NewHandler spelled in generated rpc types, even when those types are built
// from net/http ones.
package connectlike

import "connectlike/gen"

type Options struct{}

func NewHandler(opts Options, extra int) gen.CommitHandlerFunc { return nil }
