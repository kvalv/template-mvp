package object_test

import (
	"testing"

	"github.com/kvalv/template-mvp/object"
)

func TestEnvironment(t *testing.T) {
	input := struct {
		child struct {
			value int
		}
		field string
	}{
		child: struct{ value int }{
			value: 123,
		},
		field: "field",
	}

	t.Run("access field", func(t *testing.T) {
		cases := []struct {
			descr string
			key   string
			data  any
			want  object.Object
		}{
			{
				descr: "simple",
				key:   "field",
				data:  input,
				want:  &object.String{Value: "field"},
			},
			{
				descr: "nested",
				key:   "child.value",
				data:  input,
				want:  &object.Number{Value: 123},
			},
			{
				descr: "dot prefix",
				key:   ".field",
				data:  input,
				want:  &object.String{Value: "field"},
			},
			{
				descr: "dot on string",
				key:   ".",
				data:  "xx",
				want:  &object.String{Value: "xx"},
			},
			{
				descr: "dot on int",
				key:   ".",
				data:  123,
				want:  &object.Number{Value: 123},
			},
		}

		for _, tc := range cases {
			t.Run(tc.descr, func(t *testing.T) {
				env := object.NewEnvironment(tc.data)
				got := env.Field(tc.key)
				expectObjectEq(t, got, tc.want)
			})
		}
	})

	t.Run("child", func(t *testing.T) {
		got := object.NewEnvironment(input).Child("child").Field("value")
		want := &object.Number{Value: 123}
		expectObjectEq(t, got, want)
	})

	t.Run("invalid", func(t *testing.T) {
		got := object.NewEnvironment(nil).Field("field")
		if _, ok := object.AsError(got); !ok {
			t.Fatalf("expected error, got=%s", got)
		}
	})
}

func expectObjectEq(t *testing.T, got, want object.Object) {
	t.Helper()
	if got.Type() != want.Type() {
		t.Fatalf("type mismatch; want=%s, got=%s", want.Type(), got.Type())
	}
	if got.String() != want.String() {
		t.Fatalf("value mismatch; want=%s, got=%s", want.String(), got.String())
	}
}
