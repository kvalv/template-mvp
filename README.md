A re-implementation of [text/template](https://pkg.go.dev/html/template). 

The idea is to write a template engine myself, and then compare it to the way
it's implemented in the standard library.

### Current API
```go
cat := struct {
    Name     string
    LegCount int
}{
    Name:     "cat",
    LegCount: 4,
}
input := "A {{.Name}} has {{.LegCount}} legs - {{.LegCount - 2}} more than a human!"

res, err := template.New(input).Execute(&cat)
if err != nil {
    panic("oh no")
}
fmt.Println(res)
```

Outputs
```
A cat has 4 legs - 2 more than a human!
```
