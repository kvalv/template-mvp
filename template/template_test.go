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
		skip  bool
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
		{
			descr: "cond/expr/true",
			input: "{{if 2 > 1}}{{.Wow}}{{end}}",
			want:  "Wow",
			data: struct {
				Wow string
			}{Wow: "Wow"},
		},
		{
			descr: "cond/expr/false",
			input: "{{if 2 < 1}}{{.Wow}}{{end}}",
			want:  "",
			data: struct {
				Wow string
			}{Wow: "Wow"},
		},
		{
			descr: "cond/truthy/true",
			input: "{{if 2}}{{.Wow}}{{end}}",
			want:  "Wow",
			data: struct {
				Wow string
			}{Wow: "Wow"},
		},
		{
			descr: "cond/truthy/false",
			input: "{{if 0}}{{.Wow}}{{end}}",
			want:  "",
			data: struct {
				Wow string
			}{Wow: "Wow"},
		},
		{
			descr: "cond/truthy/field",
			input: "{{if .field}}X{{end}}",
			want:  "",
			data: struct {
				field int
			}{field: 0},
		},
		{
			descr: "nested",
			input: "{{if 2 > 1}}{{if 1 > 0}}hi{{end}}{{end}}",
			want:  "hi",
		},
		{
			descr: "dot",
			input: "{{.}}",
			data:  "Hello",
			want:  "Hello",
		},
		{
			descr: "range",
			input: "{{range .Slice}}Name: {{.Name}} - {{end}}",
			data: struct {
				Slice []struct {
					Name string
				}
			}{
				Slice: []struct {
					Name string
				}{
					{Name: "Alice"},
					{Name: "Bob"},
				},
			},
			skip: true,
			want: "Name: Alice - Name: Bob - ",
		},
	}

	for _, tc := range cases {
		if tc.skip {
			t.Skip(tc.descr)
		}
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
