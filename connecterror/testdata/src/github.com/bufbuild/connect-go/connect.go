package connect

type Code uint32

type Error struct{}

func (e *Error) Error() string { return "" }

func NewError(c Code, underlying error) *Error { return &Error{} }
