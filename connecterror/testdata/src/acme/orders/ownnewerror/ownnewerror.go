package ownnewerror

func NewError(err error) error { return err }

func Handler() error {
	return NewError(nil)
}
