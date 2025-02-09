package template_test

import (
	"testing"

	"github.com/kvalv/template-mvp"
)

func TestTemplate(t *testing.T) {
	cases := []struct {
		descr string
		input string
		data  any
		want  string
	}{
		{
			descr: "simple",
			input: "Hello {{.Name}}",
			data: struct {
				Name string
			}{Name: "World"},
			want: "Hello World",
		},
		{
			descr: "two variables",
			input: "One {{.Two}} {{.Three}}",
			data: struct {
				Two   string
				Three string
			}{Two: "Two", Three: "Three"},
			want: "One Two Three",
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			templ := template.New(tc.input)
			got, err := templ.Parse(&tc.data)
			if err != nil {
				t.Fatalf("Parse error: %s", err)
			}
			if tc.want != got {
				t.Fatalf("Result mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}
}
