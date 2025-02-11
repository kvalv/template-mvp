package eval_test

import (
	"testing"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/eval"
	"github.com/kvalv/template-mvp/object"
)

func TestFieldAccess(t *testing.T) {
	cases := []struct {
		descr string
		input ast.Expression
		data  any
		err   error
		want  string
	}{
		{
			descr: "struct pointer",
			input: &ast.Field{Name: "Foo"},
			data:  &struct{ Foo string }{Foo: "Bar"},
			want:  "Bar",
		},
		{
			descr: "struct value",
			input: &ast.Field{Name: "Foo"},
			data:  struct{ Foo string }{Foo: "Bar"},
			want:  "Bar",
		},
		{
			descr: "private field",
			input: &ast.Field{Name: "lower"},
			data:  struct{ lower int }{lower: 2},
			want:  "2",
		},
		{
			descr: "field not found",
			input: &ast.Field{Name: "Foo"},
			data:  struct{}{},
			err:   errors.ErrFieldNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			obj := eval.Eval(tc.input, tc.data)
			if tc.err != nil {
				expectErrorObject(t, obj, tc.err)
				return
			}
			expectNoErrorObject(t, obj)

			if got := obj.String(); tc.want != got {
				t.Fatalf("String mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}
}

func TestEvalField(t *testing.T) {
	cases := []struct {
		descr string
		expr  ast.Expression
		data  any
		want  ast.Expression
	}{
		{
			descr: "number",
			expr:  &ast.Field{Name: "myfield"},
			data:  struct{ myfield int }{myfield: 42},
			want:  &ast.Number{Value: 42},
		},
		{
			descr: "string",
			expr:  &ast.Field{Name: "myfield"},
			data:  struct{ myfield string }{myfield: "wow"},
			want:  &ast.String{Value: "wow"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			obj := eval.Eval(tc.expr, tc.data)
			if err, ok := object.AsError(obj); ok {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := obj.String(); tc.want.String() != got {
				t.Fatalf("String mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}
}

func TestEvalInfixExpression(t *testing.T) {
	cases := []struct {
		descr string
		expr  ast.Expression
		want  ast.Expression
		data  any
	}{
		{
			descr: "add numbers",
			expr: &ast.Infix{
				Lhs: &ast.Number{Value: 1},
				Rhs: &ast.Number{Value: 2},
				Op:  "+",
			},
			want: &ast.Number{Value: 3},
		},
		{
			descr: "add strings",
			expr: &ast.Infix{
				Lhs: &ast.String{Value: "foo"},
				Rhs: &ast.String{Value: "bar"},
				Op:  "+",
			},
			want: &ast.String{Value: "foobar"},
		},
		{
			descr: "number and field",
			expr: &ast.Infix{
				Lhs: &ast.Number{Value: 1},
				Rhs: &ast.Field{Name: "myfield"},
				Op:  "+",
			},
			data: struct{ myfield int }{myfield: 2},
			want: &ast.Number{Value: 3},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			obj := eval.Eval(tc.expr, tc.data)
			if err, ok := object.AsError(obj); ok {
				t.Fatalf("unexpected error: %s", err)
			}

			if got := obj.String(); tc.want.String() != got {
				t.Fatalf("String mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}

}

func expectErrorObject(t *testing.T, obj object.Object, wantErrIs error) {
	t.Helper()
	e, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("expected error object, got=%T", obj)
	}
	if !errors.Is(e, wantErrIs) {
		t.Fatalf("error mismatch; want=%q, got=%q", wantErrIs, e)
	}
}
func expectNoErrorObject(t *testing.T, obj object.Object) {
	t.Helper()
	if _, ok := obj.(*object.Error); ok {
		t.Fatalf("unexpected error: %s", obj)
	}
}
