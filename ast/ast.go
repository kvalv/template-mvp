package ast

import (
	"fmt"
	"strings"

	"github.com/kvalv/template-mvp/token"
)

type (
	Expression interface {
		String() string
	}
	Program struct {
		Exprs []Expression
	}
	Action struct {
		token.Token
		Body Expression
	}
	Text struct {
		token.Token
		Text string
	}
	Boolean struct {
		token.Token
		Value bool
	}
	Number struct {
		token.Token
		Value int
	}
	String struct {
		token.Token
		Value string
	}

	Field struct {
		token.Token
		Expression
		Name string
	}

	Cond struct {
		token.Token
		If   Expression
		Body Expression
	}

	Prefix struct {
		token.Token
		Op  string
		Rhs Expression
	}
	Infix struct {
		token.Token
		Op       string
		Lhs, Rhs Expression
	}
)

func (f *Field) String() string {
	return f.Name
}

func (n *Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (s *String) String() string {
	return s.Value
}
func (p *Program) String() string {
	var b strings.Builder
	for _, e := range p.Exprs {
		b.WriteString(e.String())
	}
	return b.String()
}
func (p *Prefix) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Rhs.String())
}
func (i *Infix) String() string {
	return fmt.Sprintf("(%s%s%s)", i.Lhs.String(), i.Op, i.Rhs.String())
}
func (a *Action) String() string {
	return fmt.Sprintf("{{%s}}", a.Body)
}
func (p *Text) String() string {
	if len(p.Text) > 50 {
		return fmt.Sprintf("%s...", p.Text[:50])
	}
	return p.Text
}
func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}
func (c *Cond) String() string {
	return fmt.Sprintf("if(%s) %s end", c.If, c.Body)
}
