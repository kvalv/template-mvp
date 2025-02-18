package eval

import (
	"reflect"

	"github.com/kvalv/template-mvp/ast"
	"github.com/kvalv/template-mvp/errors"
	"github.com/kvalv/template-mvp/object"
)

func Eval(expr ast.Expression, data any) object.Object {
	switch expr := expr.(type) {
	case *ast.Number:
		return &object.Number{Value: expr.Value}
	case *ast.String:
		return &object.String{Value: expr.Value}
	case *ast.Field:
		return evalField(expr, data)
	case *ast.Infix:
		return evalInfix(expr, data)
	case *ast.Prefix:
		return evalPrefix(expr, data)
	case *ast.Cond:
		return evalCond(expr, data)
	case *ast.Boolean:
		if expr.Value {
			return object.TRUE
		}
		return object.FALSE
	case *ast.Action:
		return Eval(expr.Body, data)
	case *ast.Text:
		return &object.String{Value: expr.Text}
	default:
		return object.Errorf("unsupported expression type %T", expr)
	}
}

func evalPrefix(expr *ast.Prefix, data any) object.Object {
	switch expr.Op {
	case ".":
		return evalField(expr.Rhs.(*ast.Field), data)
	default:
		return object.Errorf("unsupported prefix operator %s", expr.Op)
	}
}

func evalInfix(expr *ast.Infix, data any) object.Object {
	left := Eval(expr.Lhs, data)
	right := Eval(expr.Rhs, data)

	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfix(expr.Op, left.(*object.Number), right.(*object.Number))
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfix(expr.Op, left.(*object.String), right.(*object.String))

	default:
		return object.Errorf("evalInfix: unsupported types for infix expression: %v %v %v", left.Type(), expr.Op, right.Type())
	}
}

func evalNumberInfix(op string, left, right *object.Number) object.Object {
	switch op {
	case "+":
		return &object.Number{Value: left.Value + right.Value}
	case "-":
		return &object.Number{Value: left.Value - right.Value}
	case ">":
		return object.FromGoBool(left.Value > right.Value)
	case "<":
		return object.FromGoBool(left.Value < right.Value)
	case "==":
		return object.FromGoBool(left.Value == right.Value)
	default:
		return object.Errorf("unsupported operator %s", op)
	}
}

func evalStringInfix(op string, left, right *object.String) object.Object {
	switch op {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	default:
		return object.Errorf("unsupported operator %s", op)
	}
}

func evalField(expr *ast.Field, data any) object.Object {
	// we'll use reflection to access the field
	var value reflect.Value
	if v := reflect.ValueOf(data); v.Kind() == reflect.Pointer {
		value = v.Elem()
	} else {
		value = v
	}
	if data == nil {
		return object.Errorf("%w: %s", errors.ErrNilData, expr.Name)
	}
	if !value.IsValid() {
		return object.Errorf("evalField: invalid data %+v", data)
	}
	if value.Kind() != reflect.Struct {
		return object.Errorf("evalField: object is not a struct - got %+v", data)
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
func evalCond(expr *ast.Cond, data any) object.Object {
	cond := Eval(expr.If, data)
	if _, ok := object.AsError(cond); ok {
		return cond
	}
	if cond.Bool() {
		return Eval(expr.Body, data)
	}
	return &object.Void{}
}
