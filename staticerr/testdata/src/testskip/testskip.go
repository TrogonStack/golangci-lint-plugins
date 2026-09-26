// Package testskip is a fixture for confirming that staticerr rules on the
// package's own files and leaves its tests alone.
package testskip

import "errors"

var ErrNoName = errors.New("name is required")

func Validate(name string) error {
	if name == "" {
		return ErrNoName
	}

	return nil
}
