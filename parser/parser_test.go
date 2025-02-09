package parser

import (
	"testing"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/token"
)

func TestParser(t *testing.T) {
	cases := []struct {
		descr string
		input []token.Token
		want  ast.Expression
	}{
		{
			descr: "access field",
			input: []token.Token{
				{Ttype: token.DOT},
				{Ttype: token.IDENT, Text: "Name"},
			},
			want: ast.Field{
				Name: "Name",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			got, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("Parse error: %s", err)
			}
			if tc.want != got {
				t.Fatalf("Result mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}

}
