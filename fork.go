package fork

import (
	"net/http"
	"reflect"
)

type Form interface {
	New() Former
}

type Former interface {
	Form
	Fields(...Fielder) []Fielder
	Process(*http.Request)
	Valid() bool
	Render() string
}

type Field interface {
	New() Fielder
}

type Fielder interface {
	Field
	Name() string
	Get() *Value
	Set(i interface{})
	Render() string
	Valid() bool
	Errors(...string) []string
}

type Value struct {
	Raw  interface{}
	Type reflect.Type
	Kind reflect.Kind
}

func NewValue(i interface{}) *Value {
	t := reflect.TypeOf(i)
	return &Value{
		Raw:  i,
		Type: t,
		Kind: t.Kind(),
	}
}

type Widget interface {
	Render(Field) string
}
