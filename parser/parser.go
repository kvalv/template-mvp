package parser

import (
	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/token"
)

// type Parser struct {
// 	tokens []token.Token
// }

type parser struct {
	tokens []token.Token
}

func (p *parser) Parse() (ast.Expression, error) {
	if len(p.tokens) == 0 {
		return nil, errors.ErrNoTokens
	}
	// TODO: pratt parsing; we do not support .Foo + 123
	if p.tokens[0].Ttype == token.DOT {
		if len(p.tokens) != 2 {
			return nil, errors.ErrUnexpectedToken
		}
		return ast.Field{
			Name: p.tokens[1].Text,
		}, nil
	}
	return nil, errors.ErrUnexpectedToken
}

func Parse(tokens []token.Token) (ast.Expression, error) {
	p := parser{tokens: tokens}
	return p.Parse()
}
