package parser

import (
	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/token"
)

// type Parser struct {
// 	tokens []token.Token
// }

func Parse(tokens []token.Token) (ast.Expression, error) {
	return ast.Field{
		Name: "Name",
	}, nil
}

func Eval(expr *ast.Expression, v any) (string, error) {
	return "World", nil
}

// func (e *Expression) Value() string {
// 	return "World"
// }

// func (p *Parser) Parse(v any) (string, error) {
// 	return "World", nil
// }
