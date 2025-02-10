package ast

import "fmt"

type (
	Expression interface {
		String() string
	}

	Number struct {
		Value int
	}
	String struct {
		Value string
	}

	Field struct {
		Expression
		Name string
	}

	Infix struct {
		Op       string
		Lhs, Rhs Expression
	}
)

func (f Field) String() string {
	return f.Name
}

func (n Number) String() string {
	return fmt.Sprintf("%d", n.Value)
}

func (s String) String() string {
	return s.Value
}

func (i Infix) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Lhs.String(), i.Op, i.Rhs.String())
}
