package lex

import (
	"fmt"
	"io"
	"log"
	"slices"

	"github.com/kvalv/template-mvp/token"
)

type Lexer interface {
	Next() token.Token
}

type Mode int

var keywords = map[string]token.TokenType{
	"if":    token.IF,
	"end":   token.END,
	"true":  token.TRUE,
	"false": token.FALSE,
}

const (
	ModeText Mode = iota
	ModeAction
)

type lexer struct {
	log *log.Logger
	inp string
	pos int
	// whether we're inside of an action block or not
	mode Mode
	// textMode bool
}

func New(input string, logdest io.Writer) *lexer {
	log := log.New(logdest, "Lexer: ", 0)
	return &lexer{
		log: log,
		inp: input,
	}
}
func (l *lexer) curr() byte {
	if l.pos >= len(l.inp) {
		return 0
	}
	return l.inp[l.pos]
}

// retrieves the next token when the Lexer is in text mode
func (l *lexer) nextText() token.Token {
	l.log.Printf("Next(): curr=%q, peek=%q", l.curr(), l.peekNext())

	// Should we leave text mode?
	if l.curr() == '{' && l.peekNext() == '{' {
		l.log.Printf("Next(): leaving text mode, entering action mode")
		l.advance()
		l.advance()
		l.mode = ModeAction
		return token.Token{Ttype: token.ACTIONSTART, Text: "{{"}
	}

	c := l.curr()
	if c == 0 {
		return l.eof()
	}
	text := l.takewhile(func(b byte) bool {
		res := !slices.Contains([]byte{'{', '}', 0}, b)
		return res
	}, false)
	l.log.Printf("nextText(): curr=%q text=%q", c, text)
	l.advance()
	return token.Token{Ttype: token.TEXT, Text: text}

}

// retrieves the next token when the Lexer is in action mode
func (l *lexer) nextAction() token.Token {
	l.skipWhitespace()

	c := l.curr()
	switch {
	case c == 0:
		return l.eof()
	case c == '}' && l.peekNext() == '}' && l.mode == ModeAction:
		l.advance()
		l.advance()
		l.mode = ModeText
		return token.Token{Ttype: token.ACTIONEND, Text: "}}"}
	case c == '>':
		l.advance()
		return token.Token{Ttype: token.GT, Text: ">"}
	case c == '<':
		l.advance()
		return token.Token{Ttype: token.LT, Text: "<"}
	case c == '=' && l.peekNext() == '=':
		l.advance()
		l.advance()
		return token.Token{Ttype: token.EQ, Text: "=="}
	case c == '.':
		l.advance()
		return token.Token{Ttype: token.DOT, Text: "."}
	case c == '+':
		l.advance()
		return token.Token{Ttype: token.PLUS, Text: "+"}
	case c == '-':
		l.advance()
		return token.Token{Ttype: token.MINUS, Text: "-"}
	case isLetter(c):
		ident := l.takewhile(isLetter, false)
		l.advance()
		if ttype, ok := keywords[ident]; ok {
			return token.Token{Ttype: ttype, Text: ident}
		}
		return token.Token{Ttype: token.IDENT, Text: ident}
	case isDigit(c):
		num := l.takewhile(isDigit, false)
		l.advance()
		return token.Token{Ttype: token.NUMBER, Text: num}
	default:
		return l.errorf("unexpected %q", c)
	}
}

func (l *lexer) Next() token.Token {
	l.log.Printf("Next(): curr=%q", l.curr())

	if l.mode == ModeText {
		return l.nextText()
	}
	return l.nextAction()
}

func (l *lexer) takewhile(pred func(b byte) bool, consume bool) string {
	// starts at current, stops at the end
	if !pred(l.curr()) {
		return ""
	}
	start := l.pos
	for {
		c := l.peekNext()
		if c == 0 {
			break
		}
		if !pred(c) {
			if consume {
				l.advance()
			}
			break
		}
		l.advance()
	}
	end := l.pos
	return l.inp[start : end+1]
}
func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func (l *lexer) skipWhitespace() {
	for isWhitespace(l.curr()) {
		l.advance()
	}
}

func (l *lexer) peekNext() byte {
	if l.pos+1 >= len(l.inp) {
		return 0
	}
	return l.inp[l.pos+1]
}
func (l *lexer) advance() {
	l.pos++
}
func (l *lexer) slice(start, length int) string {
	return l.inp[l.pos+start : l.pos+start+length]
}

func (l *lexer) eof() token.Token {
	return token.Token{Ttype: token.EOF, Text: ""}
}

func (l *lexer) errorf(format string, a ...any) token.Token {
	l.log.Printf("Lexer.errorf: %s", fmt.Sprintf(format, a...))
	return token.Token{Ttype: token.ERROR, Text: ""}
}
