package eval

import (
	"fmt"
	"reflect"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
)

func Eval(expr ast.Expression, data any) (string, error) {
	var value reflect.Value
	if v := reflect.ValueOf(data); v.Kind() == reflect.Pointer {
		value = v.Elem()
	} else {
		value = v
	}

	if field, ok := expr.(ast.Field); ok {
		structValue := value.FieldByName(field.Name)
		if !structValue.IsValid() {
			return "", fmt.Errorf("%w: %q", errors.ErrFieldNotFound, field.Name)
		}
		return structValue.String(), nil
	}
	return "", fmt.Errorf("expr is not a field")
}
