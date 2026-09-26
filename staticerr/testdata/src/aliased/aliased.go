// Package aliased is a fixture for the two halves of resolving the callee
// through types rather than syntax: the errors package under another name is
// still the errors package, and somebody else's method called New is not.
package aliased

import goerrors "errors"

type builder struct{}

func (builder) New(string) error { return nil }

var build builder

func Aliased() error {
	return goerrors.New("name is required") // want "is built where it is returned, so every call is a different error"
}

func NotTheErrorsPackage() error {
	return build.New("name is required")
}
