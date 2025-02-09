package template

import (
	"fmt"
	"strings"

	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/eval"
	"github.com/kvalv/template-mvp/lex"
	"github.com/kvalv/template-mvp/parser"
	"github.com/kvalv/template-mvp/token"
)

type template struct {
	lexer *lex.Lexer
}

func New(input string) *template {
	return &template{
		lexer: lex.New(input),
	}
}

func (t *template) Parse(v any) (string, error) {
	out := &strings.Builder{}

	for {
		tk := t.lexer.Next()
		switch tk.Ttype {
		case token.ACTIONSTART:
			actionTokens, err := t.collectActionTokens()
			if err != nil {
				return "", err
			}
			expr, err := parser.Parse(actionTokens)
			if err != nil {
				return "", err
			}
			actionResult, err := eval.Eval(expr, v)
			if err != nil {
				return "", err
			}
			out.WriteString(actionResult)
		case token.TEXT:
			out.WriteString(tk.Text)
		case token.EOF:
			return out.String(), nil
		}
	}
}

// Consumes tokens until it finds }}, which marks the end of an action section.
// Returns an error if an unexpected token appears, e.g. EOF
func (t *template) collectActionTokens() ([]token.Token, error) {
	var res []token.Token
	for {
		tk := t.lexer.Next()
		switch tk.Ttype {
		case token.EOF, token.ACTIONSTART, token.TEXT:
			return nil, fmt.Errorf("%w: %q", errors.ErrUnexpectedToken, tk.Ttype)
		case token.ACTIONEND:
			return res, nil
		}
		res = append(res, tk)
	}
}
