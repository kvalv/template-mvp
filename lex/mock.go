package lex

import (
	"github.com/kvalv/template-mvp/token"
)

type mock struct {
	tokens []token.Token
	index  int
}

func NewMock(tk []token.Token) Lexer {
	// check if eof is present at the end. If not, add it
	if len(tk) == 0 || tk[len(tk)-1].Ttype != token.EOF {
		tk = append(tk, token.Token{Ttype: token.EOF})
	}

	return &mock{tokens: tk}
}

func (m *mock) Next() token.Token {
	i := m.index
	if i >= len(m.tokens) {
		return token.Token{Ttype: token.EOF}
	}
	m.index += 1
	return m.tokens[i]
}
