package connect

type Code uint32

const CodeInvalidArgument Code = 3

type Error struct{ code Code }

func (e *Error) Error() string { return "" }

func (e *Error) NewError() {}

func NewError(c Code, underlying error) *Error { return &Error{code: c} }
