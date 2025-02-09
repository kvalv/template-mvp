package template

type template struct {
}

func New(input string) *template {
	return &template{}
}

func (t *template) Parse(v any) (string, error) {
	return "Hello World", nil
}
