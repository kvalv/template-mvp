package ast

import (
	"fmt"

	"github.com/kvalv/template-mvp/token"
)

type (
	Expression interface {
		String() string
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

func (p *Prefix) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Rhs.String())
}
func (i *Infix) String() string {
	return fmt.Sprintf("(%s%s%s)", i.Lhs.String(), i.Op, i.Rhs.String())
}
