package object_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kvalv/template-mvp/object"
)

func TestErrorIs(t *testing.T) {
	inner := fmt.Errorf("foo")
	err := object.Errorf("%w: something bad happened", inner)

	if !errors.Is(err, inner) {
		t.Fatalf("error should not wrap inner error")
	}
}
