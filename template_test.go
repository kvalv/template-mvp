package template_test

import (
	"testing"

	"github.com/kvalv/template-mvp"
)

func TestTemplate(t *testing.T) {
	templ := template.New("Hello {{.Name}}")

	data := struct {
		Name string
	}{Name: "World"}

	got, err := templ.Parse(&data)
	if err != nil {
		t.Fatalf("Parse error: %s", err)
	}
	want := "Hello World"
	if want != got {
		t.Fatalf("Result mismatch; want=%q, got=%q", want, got)
	}
}
