package parser

import (
	"fmt"
	"io"
	"strconv"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/token"
	"github.com/kvalv/template-mvp/trace"
)

type PrefixFn func() ast.Expression
type InfixFn func(precedence int, lhs ast.Expression) ast.Expression

type parser struct {
	tr         trace.Tracer
	tokens     []token.Token
	curr, next token.Token
	index      int
	infixFns   map[token.TokenType]InfixFn
	prefixFns  map[token.TokenType]PrefixFn
}

func Parse(tokens []token.Token, logdest io.Writer) (ast.Expression, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("Parse: no tokens provided")
	}

	// artificially append EOF if not present. We want that to know where to stop
	// processing
	if hasEOF := tokens[len(tokens)-1].Ttype == token.EOF; !hasEOF {
		tokens = append(tokens, token.Token{Ttype: token.EOF})
	}

	p := parser{
		tokens:    tokens,
		tr:        trace.New(logdest),
		curr:      tokens[0],
		next:      tokens[1], // TODO: index error
		index:     1,
		infixFns:  make(map[token.TokenType]InfixFn),
		prefixFns: make(map[token.TokenType]PrefixFn),
	}

	p.prefixFns[token.IDENT] = p.parseIdentifier
	p.prefixFns[token.DOT] = p.parsePrefixExpression
	p.prefixFns[token.NUMBER] = p.parseNumber

	p.infixFns[token.PLUS] = p.parseInfixExpression
	p.infixFns[token.MINUS] = p.parseInfixExpression

	return p.Parse()
}

func (p *parser) advance() {
	// fmt.Printf("advance(): curr=%s next=%s\n", p.curr.Ttype, p.next.Ttype)
	if p.next.Ttype == token.EOF {
		p.curr = p.next
		return
	}
	p.curr = p.next
	p.index++
	p.next = p.tokens[p.index]
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	defer p.tr.Trace("parseExpression")()
	fn, ok := p.prefixFns[p.curr.Ttype]
	if !ok {
		panic(fmt.Errorf("no prefixFn found for %q", p.curr.Ttype))
	}
	expr := fn()

	// fmt.Printf("parseExpression: expr=%q, next=%q\n", expr, p.next.Ttype)

	for p.next.Ttype != token.EOF && precedence < tokenPrecedence(p.next.Ttype) {
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

func (p *parser) Parse() (expr ast.Expression, err error) {
	if len(p.tokens) == 0 {
		return nil, errors.ErrNoTokens
	}

	// on errors, the parser panics, and we catch it here.
	defer func() {
		r := recover()
		if e, ok := r.(error); ok {
			err = e
		}
	}()

	// precedence := tokenPrecedence(p.curr.Ttype)
	// fmt.Printf("precedence for %q is %d\n", p.curr.Ttype, precedence)

	expr = p.parseExpression(PrecedenceLowest)
	return expr, nil
}
