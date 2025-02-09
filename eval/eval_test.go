package eval_test

import (
	"testing"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/eval"
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
			descr: "field not found",
			input: ast.Field{Name: "Foo"},
			data:  struct{}{},
			err:   errors.ErrFieldNotFound,
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			got, err := eval.Eval(tc.input, tc.data)
			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error but got none (%q)", got)
				}
				if !errors.Is(err, tc.err) {
					t.Fatalf("error mismatch; want=%q, got=%q", tc.err, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Eval error: %s", err)
			}
			if tc.want != got {
				t.Fatalf("Result mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}

}
