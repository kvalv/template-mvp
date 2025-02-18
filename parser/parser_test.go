package parser

import (
	"os"
	"testing"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/lex"
	"github.com/kvalv/template-mvp/token"
)

func TestParser(t *testing.T) {
	cases := []struct {
		descr string
		input lex.Lexer
		want  ast.Expression
	}{
		{
			descr: "access field",
			input: lex.NewMock([]token.Token{
				{Ttype: token.DOT, Text: "."},
				{Ttype: token.IDENT, Text: "Name"},
			}),
			want: &ast.Prefix{
				Op: ".",
				Rhs: &ast.Field{
					Name: "Name",
				},
			},
		},
		{
			descr: "sum",
			input: lex.NewMock([]token.Token{
				{Ttype: token.DOT, Text: "."},
				{Ttype: token.IDENT, Text: "foo"},
				{Ttype: token.PLUS, Text: "+"},
				{Ttype: token.NUMBER, Text: "2"},
			}),
			want: &ast.Infix{
				Lhs: &ast.Prefix{
					Op: ".",
					Rhs: &ast.Field{
						Name: "foo",
					},
				},
				Op:  "+",
				Rhs: &ast.Number{Value: 2},
			},
		},
		{
			descr: "text",
			input: lex.New(`2`, os.Stderr),
			want: &ast.Text{
				Text: "2",
			},
		},
		{
			descr: "action with number",
			input: lex.New("{{2}}", os.Stderr),
			want: &ast.Action{
				Body: &ast.Number{
					Value: 2,
				},
			},
		},
		{
			descr: "cond",
			input: lex.New("{{if 1}}hi{{end}}", os.Stderr),
			want: &ast.Action{
				Body: &ast.Cond{
					If: &ast.Number{
						Value: 1,
					},
					Body: &ast.Text{
						Text: "hi",
					},
				}},
		},
		{
			descr: "greater than",
			input: lex.New("{{1 > 2}}", os.Stderr),
			want: &ast.Action{
				Body: &ast.Infix{
					Lhs: &ast.Number{
						Value: 1,
					},
					Op: ">",
					Rhs: &ast.Number{
						Value: 2,
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			parser := New(tc.input, os.Stderr)
			prog, err := parser.Parse()
			if err != nil {
				t.Fatalf("Parse error: %s", err)
			}
			if len(prog.Exprs) != 1 {
				t.Fatalf("unexpected number of expressions: %d", len(prog.Exprs))
			}
			got := prog.Exprs[0]
			t.Logf("want=%q got=%q", tc.want, got)
			expectExpression(t, tc.want, got)
		})
	}
}

func expectExpression(t *testing.T, want, got ast.Expression) {
	t.Helper()
	switch want := want.(type) {
	case *ast.Prefix:
		expectPrefix(t, want, got)
	case *ast.Number:
		expectNumber(t, want, got)
	case *ast.Field:
		expectField(t, want, got)
	case *ast.Infix:
		expectInfix(t, want, got)
	case *ast.Text:
		expectText(t, want, got)
	case *ast.Action:
		expectAction(t, want, got)
	case *ast.Cond:
		expectCond(t, want, got)
	default:
		t.Fatalf("unexpected type: %T", want)
	}
}

func expectPrefix(t *testing.T, want *ast.Prefix, got ast.Expression) {
	t.Helper()
	prefix, ok := got.(*ast.Prefix)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	if prefix.Op != want.Op {
		t.Fatalf("op mismatch; want=%q, got=%q", want.Op, prefix.Op)
	}
	expectExpression(t, want.Rhs, prefix.Rhs)
}

func expectAction(t *testing.T, want *ast.Action, got ast.Expression) {
	t.Helper()
	action, ok := got.(*ast.Action)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	expectExpression(t, want.Body, action.Body)
}
func expectCond(t *testing.T, want *ast.Cond, got ast.Expression) {
	t.Helper()
	cond, ok := got.(*ast.Cond)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	expectExpression(t, want.If, cond.If)
	expectExpression(t, want.Body, cond.Body)
}

func expectNumber(t *testing.T, want *ast.Number, got ast.Expression) {
	t.Helper()
	number, ok := got.(*ast.Number)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	if number.Value != want.Value {
		t.Fatalf("value mismatch; want=%d, got=%d", want.Value, number.Value)
	}
}
func expectField(t *testing.T, want *ast.Field, got ast.Expression) {
	t.Helper()
	field, ok := got.(*ast.Field)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	if field.Name != want.Name {
		t.Fatalf("name mismatch; want=%q, got=%q", want.Name, field.Name)
	}
}
func expectInfix(t *testing.T, want *ast.Infix, got ast.Expression) {
	t.Helper()
	infix, ok := got.(*ast.Infix)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	if infix.Op != want.Op {
		t.Fatalf("op mismatch; want=%q, got=%q", want.Op, infix.Op)
	}
	expectExpression(t, want.Lhs, infix.Lhs)
	expectExpression(t, want.Rhs, infix.Rhs)
}
func expectText(t *testing.T, want *ast.Text, got ast.Expression) {
	t.Helper()
	text, ok := got.(*ast.Text)
	if !ok {
		t.Fatalf("type mismatch; want=%T, got=%T", want, got)
	}
	if want.Text != text.Text {
		t.Fatalf("text mismatch; want=%q, got=%q", want.Text, text.Text)
	}
}
