package template

import "github.com/kvalv/template-mvp/token"

// type Parser struct {
// 	tokens []token.Token
// }

func Parse(tokens []token.Token) (*Expression, error) {
	return &Expression{}, nil
}

type Expression struct {
}

func Eval(expr *Expression, v any) (string, error) {
	return "World", nil
}

// func (e *Expression) Value() string {
// 	return "World"
// }

// func (p *Parser) Parse(v any) (string, error) {
// 	return "World", nil
// }
