// Package localvar is a fixture for confirming that a var declaration inside
// a function buys no exemption: what the rule is about is where the value is
// built, not the keyword it is built under.
package localvar

import "errors"

func Validate(name string) error {
	var errNoName = errors.New("name is required") // want "is built where it is returned, so every call is a different error"

	if name == "" {
		return errNoName
	}

	return nil
}
