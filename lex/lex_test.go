package lex_test

import (
	"os"
	"testing"

	"github.com/kvalv/template-mvp/lex"
	"github.com/kvalv/template-mvp/token"
)

func TestLexer(t *testing.T) {
	cases := []struct {
		descr string
		input string
		want  []token.Token
	}{
		{
			descr: "text",
			input: "hello",
			want: []token.Token{
				{Ttype: token.TEXT, Text: "hello"},
				{Ttype: token.EOF, Text: ""},
			},
		},
		{
			descr: "action",
			input: "{{.World}}",
			want: []token.Token{
				{Ttype: token.ACTIONSTART, Text: "{{"},
				{Ttype: token.DOT, Text: "."},
				{Ttype: token.IDENT, Text: "World"},
				{Ttype: token.ACTIONEND, Text: "}}"},
				{Ttype: token.EOF, Text: ""},
			},
		},
		{
			descr: "whitespace",
			input: "{{ World   }}",
			want: []token.Token{
				{Ttype: token.ACTIONSTART, Text: "{{"},
				{Ttype: token.IDENT, Text: "World"},
				{Ttype: token.ACTIONEND, Text: "}}"},
				{Ttype: token.EOF, Text: ""},
			},
		},
		{
			descr: "arithmetic",
			input: "{{2 + 3}}",
			want: []token.Token{
				{Ttype: token.ACTIONSTART, Text: "{{"},
				{Ttype: token.NUMBER, Text: "2"},
				{Ttype: token.PLUS, Text: "+"},
				{Ttype: token.NUMBER, Text: "3"},
				{Ttype: token.ACTIONEND, Text: "}}"},
			},
		},
		{
			descr: "if block",
			input: "{{if true}}Hello{{end}}",
			want: []token.Token{
				{Ttype: token.ACTIONSTART, Text: "{{"},
				{Ttype: token.IF, Text: "if"},
				{Ttype: token.TRUE, Text: "true"},
				{Ttype: token.ACTIONEND, Text: "}}"},
				{Ttype: token.TEXT, Text: "Hello"},
				{Ttype: token.ACTIONSTART, Text: "{{"},
				{Ttype: token.END, Text: "end"},
				{Ttype: token.ACTIONEND, Text: "}}"},
				{Ttype: token.EOF, Text: ""},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.descr, func(t *testing.T) {
			lexer := lex.New(tc.input, os.Stderr)
			for _, tk := range tc.want {
				got := lexer.Next()
				expectTokenMatch(t, got, tk)
			}
		})
	}
}

func expectTokenMatch(t *testing.T, got, want token.Token) {
	// TODO: spans?
	t.Helper()
	if got.Ttype != want.Ttype {
		t.Fatalf("TokenType mismatch: got=%q, want=%q (Text=%q)", got.Ttype, want.Ttype, got.Text)
	}
	if got.Text != want.Text {
		t.Fatalf("Text mismatch: got=%q, want=%q", got.Text, want.Text)
	}
}
