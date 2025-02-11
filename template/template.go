package template

import (
	"fmt"
	"io"
	"strings"

	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/eval"
	"github.com/kvalv/template-mvp/lex"
	"github.com/kvalv/template-mvp/object"
	"github.com/kvalv/template-mvp/parser"
	"github.com/kvalv/template-mvp/token"
)

type template struct {
	logdest io.Writer
	lexer   *lex.Lexer
}

func New(input string, logdest io.Writer) *template {
	return &template{
		logdest: logdest,
		lexer:   lex.New(input, logdest),
	}
}

func (t *template) debugTokens(tks []token.Token) {
	b := strings.Builder{}
	for _, tk := range tks {
		b.WriteString(string(tk.Ttype))
		b.WriteString(" ")
	}
	fmt.Fprintf(t.logdest, "debugTokens: %s\n", b.String())
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
			t.debugTokens(actionTokens)
			expr, err := parser.Parse(actionTokens, t.logdest)
			if err != nil {
				return "", err
			}
			actionResult := eval.Eval(expr, v)
			if err, ok := object.AsError(actionResult); ok {
				return "", err
			}
			fmt.Fprintf(out, "%s", actionResult)
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
