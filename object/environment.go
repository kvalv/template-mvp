package object

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kvalv/template-mvp/errors"
)

type Environment struct {
	data reflect.Value
}

func NewEnvironment(input any) *Environment {
	var data reflect.Value
	if v, ok := input.(reflect.Value); ok {
		data = v
	} else {
		data = reflect.ValueOf(input)
	}

	return &Environment{
		data: data,
	}
}

func (e *Environment) field(name string) (reflect.Value, error) {
	structValue := e.data
	for _, part := range strings.Split(name, ".") {
		if part == "" {
			continue
		}
		if !structValue.IsValid() {
			return reflect.Value{}, fmt.Errorf("%w: %s", errors.ErrFieldNotFound, name)
		}
		structValue = structValue.FieldByName(part)
		if !structValue.IsValid() {
			return reflect.Value{}, fmt.Errorf("%w: %s", errors.ErrFieldNotFound, name)
		}
	}
	return structValue, nil
}

// Field returns the value of the field at the given path.
func (e *Environment) Field(path string) Object {
	value, err := e.field(path)
	if err != nil {
		return Errorf(err.Error())
	}

	switch value.Kind() {
	case reflect.String:
		return &String{Value: value.String()}
	case reflect.Int:
		return &Number{Value: int(value.Int())}
	default:
		return Errorf("AccessField: field %q - unsupported type %s ", path, value.Kind())
	}
}

// Child returns a new environment with the field at the given path.
func (e *Environment) Child(path string) *Environment {
	field, err := e.field(path)
	if err != nil {
		return &Environment{
			data: reflect.Value{},
		}
	}
	return &Environment{
		data: field,
	}
}
