package eval

import (
	"reflect"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/object"
)

func Eval(expr ast.Expression, data any) object.Object {
	switch expr := expr.(type) {
	case ast.Field:
		return evalField(expr, data)
	default:
		return object.Errorf("unsupported expression type %T", expr)
	}
}

func evalField(expr ast.Field, data any) object.Object {
	// we'll use reflection to access the field
	var value reflect.Value
	if v := reflect.ValueOf(data); v.Kind() == reflect.Pointer {
		value = v.Elem()
	} else {
		value = v
	}

	structValue := value.FieldByName(expr.Name)
	if !structValue.IsValid() {
		return object.Errorf("%w: %s", errors.ErrFieldNotFound, expr.Name)
	}

	switch structValue.Kind() {
	case reflect.String:
		return &object.String{Value: structValue.String()}
	case reflect.Int:
		return &object.Number{Value: int(structValue.Int())}
	default:
		return object.Errorf("unsupported type %s", structValue.Kind())
	}
}
