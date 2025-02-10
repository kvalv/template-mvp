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
			input: ast.Field{Name: "Foo"},
			data:  &struct{ Foo string }{Foo: "Bar"},
			want:  "Bar",
		},
		{
			descr: "struct value",
			input: ast.Field{Name: "Foo"},
			data:  struct{ Foo string }{Foo: "Bar"},
			want:  "Bar",
		},
		{
			descr: "private field",
			input: ast.Field{Name: "lower"},
			data:  struct{ lower int }{lower: 2},
			want:  "2",
		},
		{
			descr: "field not found",
			input: ast.Field{Name: "Foo"},
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
