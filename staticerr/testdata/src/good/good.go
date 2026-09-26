// Package good is a fixture for confirming that staticerr reports nothing
// when every errors.New is a package-level var, in each of the shapes one is
// written in.
package good

import "errors"

// A var on its own.
var ErrAlone = errors.New("alone")

// An unexported var, which is still declared once.
var errHidden = errors.New("hidden")

// A var block, which is how a package with more than one sentinel declares
// them.
var (
	ErrFirst  = errors.New("first")
	ErrSecond = errors.New("second")
)

// Inside a composite literal, still at package level and still built once.
var byCode = map[int]error{
	404: errors.New("not found"),
	409: errors.New("already exists"),
}

func Lookup(code int) error {
	if err, ok := byCode[code]; ok {
		return err
	}

	if code == 0 {
		return errHidden
	}

	return ErrAlone
}
