// Package inline is a fixture for confirming that staticerr reports an
// errors.New built where it is returned.
package inline

import "errors"

func Validate(name string) error {
	if name == "" {
		return errors.New("name is required") // want "is built where it is returned, so every call is a different error"
	}

	return nil
}
