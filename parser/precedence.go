package parser

import (
	"github.com/kvalv/template-mvp/token"
)

const (
	_ int = iota
	PrecedenceLowest
	PrecedencePlus
	PrecedenceMul
	PrecedencePrefix
)

func tokenPrecedence(ttype token.TokenType) (p int) {
	switch ttype {
	case token.DOT:
		return PrecedencePrefix
	case token.PLUS, token.MINUS, token.GT, token.LT, token.EQ:
		return PrecedencePlus
	default:
		return PrecedenceLowest
	}
}
