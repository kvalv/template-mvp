package token

type TokenType string

const (
	ERROR       TokenType = "ERROR"
	EOF         TokenType = "EOF"
	TEXT        TokenType = "TEXT"
	ACTIONSTART TokenType = "ACTIONSTART"
	ACTIONEND   TokenType = "ACTIONEND"
	DOT         TokenType = "DOT"
	IDENT       TokenType = "IDENT"
	NUMBER      TokenType = "NUMBER"
	PLUS        TokenType = "PLUS"
)

type Token struct {
	Ttype TokenType
	Text  string
	Span
}

type Span struct {
	Start, End Position
}
type Position struct {
	Row, Col int
}
