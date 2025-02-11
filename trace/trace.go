package trace

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Tracer interface {
	Trace(string) func()
}

type tracer struct {
	log.Logger
	level int
}

func New(w io.Writer) *tracer {
	return &tracer{
		Logger: *log.New(w, "", 0),
	}
}
func (t *tracer) Trace(name string) func() {
	start := time.Now()
	indent := strings.Repeat(" ", t.level*2)
	t.level++
	t.Logger.Printf("%sBEGIN %s", indent, name)
	return func() {
		elapsed := time.Since(start)
		line := fmt.Sprintf("%sEND %s [%s]", indent, name, elapsed.String())
		t.Logger.Printf(line)
		t.level--
	}
}
