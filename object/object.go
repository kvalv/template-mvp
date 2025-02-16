package object

import "fmt"

type ObjectType string

const (
	STRING_OBJ  = "STRING"
	NUMBER_OBJ  = "NUMBER"
	ERROR_OBJ   = "ERROR"
	BOOLEAN_OBJ = "BOOLEAN"
)

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	String() string
}

type (
	String  struct{ Value string }
	Number  struct{ Value int }
	Error   struct{ err error }
	Boolean struct{ Value bool }
)

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) String() string   { return s.Value }

func (n *Number) Type() ObjectType { return NUMBER_OBJ }
func (n *Number) String() string   { return fmt.Sprintf("%d", n.Value) }

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) String() string   { return fmt.Sprintf("%t", b.Value) }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) String() string   { return e.err.Error() }
func (e *Error) Unwrap() error    { return e.err }
func (e *Error) Error() string    { return e.err.Error() }

func Errorf(format string, args ...interface{}) *Error {
	return &Error{err: fmt.Errorf(format, args...)}
}
func AsError(obj Object) (error, bool) {
	err, ok := obj.(*Error)
	return err, ok
}
