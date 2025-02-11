package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kvalv/template-mvp/template"
)

func main() {
	cat := struct {
		Name     string
		LegCount int
	}{
		Name:     "cat",
		LegCount: 4,
	}
	input := "A {{.Name}} has {{.LegCount }} legs - {{.LegCount - 2}} more than a human!"
	res, err := template.New(input, io.Discard).Parse(&cat)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
	}
	fmt.Println(res)
}
