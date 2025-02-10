package errors

import "errors"

// re-export
var (
	New = errors.New
	Is  = errors.Is
	As  = errors.As
)

var (
	ErrUnexpectedToken = errors.New("unexpected token")
	ErrNoTokens        = errors.New("no tokens")
	ErrFieldNotFound   = errors.New("Field not found")
)
