// Package nested is a fixture for an errors.New handed straight to something
// else, such as a constructor that pairs a status code with a message in one
// expression.
package nested

import (
	"errors"
	"fmt"
)

func Wrapped() error {
	return fmt.Errorf("validating: %w", errors.New("name is required")) // want "is built where it is returned, so every call is a different error"
}
