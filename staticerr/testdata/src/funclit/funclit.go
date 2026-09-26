// Package funclit is a fixture for the case a naive search would get wrong: a
// package-level var holding a func literal. The var declaration is above the
// literal, but the errors.New inside it runs on every call of the func and not
// once at initialisation, so it is reported like any other.
package funclit

import "errors"

var Validate = func(name string) error {
	if name == "" {
		return errors.New("name is required") // want "is built where it is returned, so every call is a different error"
	}

	return nil
}
