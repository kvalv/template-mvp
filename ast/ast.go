package ast

type (
	Expression interface {
		String() string
	}

	Field struct {
		Expression
		Name string
	}
)

func (f Field) String() string {
	return f.Name
}
