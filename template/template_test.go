package template_test

import (
	"os"
	"testing"

	"github.com/kvalv/template-mvp/template"
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
		{
			descr: "sum",
			input: "{{.One + 2}}",
			data: struct {
				One int
			}{One: 1},
			want: "3",
		},
		{
			descr: "cond/true",
			input: "{{if true}}hi{{end}}",
			want:  "hi",
			data:  nil,
		},
		{
			descr: "cond/false",
			input: "{{if false}}hi{{end}}",
			want:  "",
			data:  nil,
		},
		{
			descr: "cond/variable",
			input: "{{if true}}{{.Wow}}{{end}}",
			want:  "Wow",
			data: struct {
				Wow string
			}{Wow: "Wow"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			templ := template.New(tc.input, template.LogDest(os.Stderr))
			got, err := templ.Execute(tc.data)
			if err != nil {
				t.Fatalf("Parse error: %s", err)
			}
			if tc.want != got {
				t.Fatalf("Result mismatch; want=%q, got=%q", tc.want, got)
			}
		})
	}
}
