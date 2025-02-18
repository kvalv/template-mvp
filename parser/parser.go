package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/lex"
	"github.com/kvalv/template-mvp/token"
	"github.com/kvalv/template-mvp/trace"
)

type PrefixFn func() ast.Expression
type InfixFn func(precedence int, lhs ast.Expression) ast.Expression

type parser struct {
	lex        lex.Lexer
	curr, next token.Token

	infixFns  map[token.TokenType]InfixFn
	prefixFns map[token.TokenType]PrefixFn

	tr trace.Tracer
}

func New(lex lex.Lexer, logdest io.Writer) *parser {
	p := &parser{
		prefixFns: make(map[token.TokenType]PrefixFn),
		infixFns:  make(map[token.TokenType]InfixFn),
		tr:        trace.New(logdest),
		lex:       lex,
	}

	p.prefixFns[token.ACTIONSTART] = p.parseAction
	p.prefixFns[token.TEXT] = p.parseText
	p.prefixFns[token.IDENT] = p.parseIdentifier
	p.prefixFns[token.DOT] = p.parsePrefixExpression
	p.prefixFns[token.NUMBER] = p.parseNumber
	p.prefixFns[token.IF] = p.parseCond
	p.prefixFns[token.TRUE] = p.parseBoolean
	p.prefixFns[token.FALSE] = p.parseBoolean

	for _, tk := range []token.TokenType{token.PLUS, token.MINUS, token.GT, token.LT, token.EQ} {
		p.infixFns[tk] = p.parseInfixExpression
	}

	p.advance()
	p.advance()
	return p

}

func (p *parser) advance() {
	p.curr = p.next
	if p.next.Ttype == token.EOF {
		return
	}
	p.next = p.lex.Next()
}

func (p *parser) parseText() ast.Expression {
	if p.curr.Ttype != token.TEXT {
		panic(fmt.Errorf("expected token type TEXT, got %s", p.curr.Ttype))
	}

	defer p.tr.Trace("parseText")()
	return &ast.Text{
		Text:  p.curr.Text,
		Token: p.curr,
	}
}

func (p *parser) parseAction() ast.Expression {
	// an action is delimited by {{ and }}
	defer p.tr.Trace("parseAction")
	p.expectToken(token.ACTIONSTART)
	expr := &ast.Action{
		Token: p.curr,
	}
	p.advance()
	expr.Body = p.parseExpression(PrecedenceLowest)

	// The last token for Conditionals is ACTIONEND, but for other
	// expressions it's not. So we just ensure that in any case,
	// The last token is ACTIONEND. 🤷
	if p.next.Ttype == token.ACTIONEND {
		p.advance()
	}

	p.expectToken(token.ACTIONEND)
	return expr
}

func (p *parser) parseCond() ast.Expression {
	defer p.tr.Trace("parseCond")()
	p.expectToken(token.IF)
	cond := &ast.Cond{
		Token: p.curr,
	}
	p.advance()
	cond.If = p.parseExpression(PrecedenceLowest)
	p.advance()
	p.expectToken(token.ACTIONEND)
	p.advance()
	cond.Body = p.parseExpression(PrecedenceLowest)
	p.advance()

	p.expectToken(token.ACTIONSTART)
	p.advance()
	p.expectToken(token.END)
	p.advance()
	p.expectToken(token.ACTIONEND)

	return cond
}

func (p *parser) parseBoolean() ast.Expression {
	defer p.tr.Trace("parseBoolean")()
	return &ast.Boolean{
		Token: p.curr,
		Value: p.curr.Ttype == token.TRUE,
	}
}

func (p *parser) expectToken(ttype token.TokenType, extra ...string) {
	if p.curr.Ttype != ttype {
		msg := fmt.Sprintf("expected token type %q, got %q", ttype, p.curr.Ttype)
		for _, e := range extra {
			msg += " " + e
		}
		panic(msg)
	}
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	defer p.tr.Trace("parseExpression")()
	fn, ok := p.prefixFns[p.curr.Ttype]
	if !ok {
		panic(fmt.Errorf("no prefixFn found for %q", p.curr.Ttype))
	}
	expr := fn()

	// fmt.Printf("parseExpression: expr=%q, next=%q\n", expr, p.next.Ttype)

	for p.next.Ttype != token.EOF && p.next.Ttype != token.ACTIONEND && precedence < tokenPrecedence(p.next.Ttype) {
		p.advance()
		infixFn, ok := p.infixFns[p.curr.Ttype]
		if !ok {
			panic(fmt.Errorf("no infixFn found for %q", p.curr.Ttype))
		}
		expr = infixFn(
			tokenPrecedence(p.curr.Ttype),
			expr,
		)
	}

	return expr
}

func (p *parser) parsePrefixExpression() ast.Expression {
	defer p.tr.Trace("parsePrefixExpression")()
	// current is dot, next is then a field
	expr := &ast.Prefix{
		Token: p.curr,
		Op:    p.curr.Text,
	}
	p.advance()
	if expr.Rhs = p.parseExpression(PrecedencePrefix); expr.Rhs == nil {
		return nil
	}
	return expr
}

func (p *parser) parseInfixExpression(precedence int, lhs ast.Expression) ast.Expression {
	defer p.tr.Trace("parseInfixExpression")()
	expr := &ast.Infix{
		Token: p.curr,
		Lhs:   lhs,
		Op:    p.curr.Text,
	}
	p.advance()
	expr.Rhs = p.parseExpression(precedence)
	return expr
}

func (p *parser) parseIdentifier() ast.Expression {
	defer p.tr.Trace("parseIdentifier")()
	return &ast.Field{
		Token: p.curr,
		Name:  p.curr.Text,
	}
}

func (p *parser) parseNumber() ast.Expression {
	defer p.tr.Trace("parseNumber")()
	value, err := strconv.Atoi(p.curr.Text)
	if err != nil {
		panic(fmt.Errorf("parseNumber: not a number: %q", p.curr.Text))
	}
	return &ast.Number{
		Token: p.curr,
		Value: value,
	}
}

// Parse returns the next expression
func (p *parser) Parse() (prog *ast.Program, err error) {
	if p.curr.Ttype == token.EOF {
		return nil, errors.ErrNoTokens
	}

	// on errors, the parser panics, and we catch it here.
	defer func() {
		r := recover()
		if e, ok := r.(error); ok {
			err = e
		} else if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	// precedence := tokenPrecedence(p.curr.Ttype)
	// fmt.Printf("precedence for %q is %d\n", p.curr.Ttype, precedence)

	var res ast.Program
	var it int
	for p.curr.Ttype != token.EOF {
		expr := p.parseExpression(PrecedenceLowest)
		p.advance()
		res.Exprs = append(res.Exprs, expr)
		it++
		if it > 5 {
			break
		}
	}
	prog = &res
	return
}
